/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	korifiv1alpha1 "code.cloudfoundry.org/korifi/controllers/api/v1alpha1"
	"code.cloudfoundry.org/korifi/controllers/cleanup"
	"code.cloudfoundry.org/korifi/controllers/config"
	networkingcontrollers "code.cloudfoundry.org/korifi/controllers/controllers/networking"
	servicescontrollers "code.cloudfoundry.org/korifi/controllers/controllers/services"
	"code.cloudfoundry.org/korifi/controllers/controllers/shared"
	workloadscontrollers "code.cloudfoundry.org/korifi/controllers/controllers/workloads"
	"code.cloudfoundry.org/korifi/controllers/controllers/workloads/env"
	"code.cloudfoundry.org/korifi/controllers/controllers/workloads/labels"
	"code.cloudfoundry.org/korifi/controllers/coordination"
	"code.cloudfoundry.org/korifi/controllers/webhooks"
	controllersfinalizer "code.cloudfoundry.org/korifi/controllers/webhooks/finalizer"
	"code.cloudfoundry.org/korifi/controllers/webhooks/networking"
	"code.cloudfoundry.org/korifi/controllers/webhooks/services"
	versionwebhook "code.cloudfoundry.org/korifi/controllers/webhooks/version"
	"code.cloudfoundry.org/korifi/controllers/webhooks/workloads"
	jobtaskrunnercontrollers "code.cloudfoundry.org/korifi/job-task-runner/controllers"
	"code.cloudfoundry.org/korifi/kpack-image-builder/controllers"
	kpackimagebuilderfinalizer "code.cloudfoundry.org/korifi/kpack-image-builder/controllers/webhooks/finalizer"
	statesetfulrunnerv1 "code.cloudfoundry.org/korifi/statefulset-runner/api/v1"
	statefulsetcontrollers "code.cloudfoundry.org/korifi/statefulset-runner/controllers"
	"code.cloudfoundry.org/korifi/tools"
	"code.cloudfoundry.org/korifi/tools/image"
	"code.cloudfoundry.org/korifi/tools/registry"
	"code.cloudfoundry.org/korifi/version"

	buildv1alpha2 "github.com/pivotal/kpack/pkg/apis/build/v1alpha2"
	contourv1 "github.com/projectcontour/contour/apis/projectcontour/v1"
	servicebindingv1beta1 "github.com/servicebinding/runtime/apis/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	k8sclient "k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/v2"
	admission "k8s.io/pod-security-admission/api"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(buildv1alpha2.AddToScheme(scheme))
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(contourv1.AddToScheme(scheme))
	utilruntime.Must(korifiv1alpha1.AddToScheme(scheme))
	utilruntime.Must(servicebindingv1beta1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate -o fake -fake-name Client sigs.k8s.io/controller-runtime/pkg/client.Client
//counterfeiter:generate -o fake -fake-name EventRecorder k8s.io/client-go/tools/record.EventRecorder
//counterfeiter:generate -o fake -fake-name StatusWriter sigs.k8s.io/controller-runtime/pkg/client.StatusWriter

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string

	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

	configPath, found := os.LookupEnv("CONTROLLERSCONFIG")
	if !found {
		panic("CONTROLLERSCONFIG must be set")
	}

	controllerConfig, err := config.LoadFromPath(configPath)
	if err != nil {
		errorMessage := fmt.Sprintf("Config could not be read: %v", err)
		panic(errorMessage)
	}

	logger, atomicLevel, err := tools.NewZapLogger(controllerConfig.LogLevel)
	if err != nil {
		panic(fmt.Sprintf("error creating new zap logger: %v", err))
	}

	ctrl.SetLogger(logger)
	klog.SetLogger(ctrl.Log)

	log.SetOutput(&tools.LogrWriter{Logger: ctrl.Log, Message: "HTTP server error"})

	ctrl.Log.Info("starting Korifi controllers", "version", version.Version)

	conf := ctrl.GetConfigOrDie()
	k8sClient, err := k8sclient.NewForConfig(conf)
	if err != nil {
		panic(fmt.Sprintf("could not create k8s client: %v", err))
	}

	mgr, err := ctrl.NewManager(conf, ctrl.Options{
		Scheme: scheme,
		WebhookServer: webhook.NewServer(webhook.Options{
			Port: 9443,
		}),
		Metrics: metricsserver.Options{
			BindAddress: metricsAddr,
		},
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "13c200ec.cloudfoundry.org",
	})
	if err != nil {
		setupLog.Error(err, "unable to initialize manager")
		os.Exit(1)
	}

	if os.Getenv("ENABLE_CONTROLLERS") != "false" {
		imageClient := image.NewClient(k8sClient)

		if err = (workloadscontrollers.NewCFAppReconciler(
			mgr.GetClient(),
			mgr.GetScheme(),
			ctrl.Log.WithName("controllers").WithName("CFApp"),
			env.NewVCAPServicesEnvValueBuilder(mgr.GetClient()),
			env.NewVCAPApplicationEnvValueBuilder(mgr.GetClient(), controllerConfig.ExtraVCAPApplicationValues),
		)).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "CFApp")
			os.Exit(1)
		}

		buildCleaner := cleanup.NewBuildCleaner(mgr.GetClient(), controllerConfig.MaxRetainedBuildsPerApp)
		if err = (workloadscontrollers.NewCFBuildpackBuildReconciler(
			mgr.GetClient(),
			buildCleaner,
			mgr.GetScheme(),
			ctrl.Log.WithName("controllers").WithName("CFBuildpackBuild"),
			controllerConfig,
			env.NewWorkloadEnvBuilder(mgr.GetClient()),
		)).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "CFBuildpackBuild")
			os.Exit(1)
		}

		if err = (workloadscontrollers.NewCFDockerBuildReconciler(
			mgr.GetClient(),
			buildCleaner,
			imageClient,
			mgr.GetScheme(),
			ctrl.Log.WithName("controllers").WithName("CFDockerBuild"),
		)).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "CFDockerBuild")
			os.Exit(1)
		}

		if err = (workloadscontrollers.NewCFPackageReconciler(
			mgr.GetClient(),
			mgr.GetScheme(),
			ctrl.Log.WithName("controllers").WithName("CFPackage"),
			imageClient,
			cleanup.NewPackageCleaner(mgr.GetClient(), controllerConfig.MaxRetainedPackagesPerApp),
			controllerConfig.ContainerRegistrySecretNames,
		)).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "CFPackage")
			os.Exit(1)
		}

		if err = (workloadscontrollers.NewCFProcessReconciler(
			mgr.GetClient(),
			mgr.GetScheme(),
			ctrl.Log.WithName("controllers").WithName("CFProcess"),
			controllerConfig,
			env.NewWorkloadEnvBuilder(mgr.GetClient()),
		)).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "CFProcess")
			os.Exit(1)
		}

		if err = (servicescontrollers.NewCFServiceInstanceReconciler(
			mgr.GetClient(),
			mgr.GetScheme(),
			ctrl.Log.WithName("controllers").WithName("CFServiceInstance"),
		)).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "CFServiceInstance")
			os.Exit(1)
		}

		if err = (servicescontrollers.NewCFServiceBindingReconciler(
			mgr.GetClient(),
			mgr.GetScheme(),
			ctrl.Log.WithName("controllers").WithName("CFServiceBinding"),
		)).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "CFServiceBinding")
			os.Exit(1)
		}

		labelCompiler := labels.NewCompiler().
			Defaults(map[string]string{
				admission.EnforceLevelLabel: string(admission.LevelRestricted),
				admission.AuditLevelLabel:   string(admission.LevelRestricted),
			}).
			Defaults(controllerConfig.NamespaceLabels)

		if err = workloadscontrollers.NewCFOrgReconciler(
			mgr.GetClient(),
			ctrl.Log.WithName("controllers").WithName("CFOrg"),
			controllerConfig.ContainerRegistrySecretNames,
			labelCompiler,
		).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "CFOrg")
			os.Exit(1)
		}

		if err = workloadscontrollers.NewCFSpaceReconciler(
			mgr.GetClient(),
			ctrl.Log.WithName("controllers").WithName("CFSpace"),
			controllerConfig.ContainerRegistrySecretNames,
			controllerConfig.CFRootNamespace,
			*controllerConfig.SpaceFinalizerAppDeletionTimeout,
			labelCompiler,
		).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "CFSpace")
			os.Exit(1)
		}

		var taskTTL time.Duration
		taskTTL, err = controllerConfig.ParseTaskTTL()
		if err != nil {
			setupLog.Error(err, "failed to parse task TTL", "controller", "CFTask", "taskTTL", controllerConfig.TaskTTL)
			os.Exit(1)

		}
		if err = workloadscontrollers.NewCFTaskReconciler(
			mgr.GetClient(),
			mgr.GetScheme(),
			mgr.GetEventRecorderFor("cftask-controller"),
			ctrl.Log.WithName("controllers").WithName("CFTask"),
			env.NewWorkloadEnvBuilder(mgr.GetClient()),
			taskTTL,
		).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "CFTask")
			os.Exit(1)
		}

		if err = (networkingcontrollers.NewCFDomainReconciler(
			mgr.GetClient(),
			mgr.GetScheme(),
			ctrl.Log.WithName("controllers").WithName("CFDomain"),
		)).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "CFDomain")
			os.Exit(1)
		}
		//+kubebuilder:scaffold:builder

		// Setup Index with Manager
		err = shared.SetupIndexWithManager(mgr)
		if err != nil {
			setupLog.Error(err, "unable to setup index on manager")
			os.Exit(1)
		}

		if controllerConfig.IncludeKpackImageBuilder {
			var builderReadinessTimeout time.Duration
			builderReadinessTimeout, err = controllerConfig.ParseBuilderReadinessTimeout()
			if err != nil {
				setupLog.Error(err, "error parsing builderReadinessTimeout")
				os.Exit(1)
			}
			if err = controllers.NewBuildWorkloadReconciler(
				mgr.GetClient(),
				mgr.GetScheme(),
				ctrl.Log.WithName("controllers").WithName("BuildWorkloadReconciler"),
				controllerConfig,
				imageClient,
				controllerConfig.ContainerRepositoryPrefix,
				registry.NewRepositoryCreator(controllerConfig.ContainerRegistryType),
				builderReadinessTimeout,
			).SetupWithManager(mgr); err != nil {
				setupLog.Error(err, "unable to create controller", "controller", "BuildWorkload")
				os.Exit(1)
			}

			if err = controllers.NewBuilderInfoReconciler(
				mgr.GetClient(),
				mgr.GetScheme(),
				ctrl.Log.WithName("controllers").WithName("BuilderInfoReconciler"),
				controllerConfig.ClusterBuilderName,
				controllerConfig.CFRootNamespace,
			).SetupWithManager(mgr); err != nil {
				setupLog.Error(err, "unable to create controller", "controller", "BuilderInfo")
				os.Exit(1)
			}

			if err = controllers.NewKpackBuildController(
				mgr.GetClient(),
				ctrl.Log.WithName("kpack-image-builder").WithName("KpackBuild"),
				imageClient,
				controllerConfig.BuilderServiceAccount,
			).SetupWithManager(mgr); err != nil {
				setupLog.Error(err, "unable to create controller", "controller", "KpackBuild")
				os.Exit(1)
			}
		}

		if controllerConfig.IncludeJobTaskRunner {
			logger := ctrl.Log.WithName("controllers").WithName("TaskWorkload")
			var jobTTL time.Duration
			jobTTL, err = controllerConfig.ParseJobTTL()
			if err != nil {
				panic(err)
			}

			taskWorkloadReconciler := jobtaskrunnercontrollers.NewTaskWorkloadReconciler(
				logger,
				mgr.GetClient(),
				mgr.GetScheme(),
				jobtaskrunnercontrollers.NewStatusGetter(logger, mgr.GetClient()),
				jobTTL,
			)
			if err = taskWorkloadReconciler.SetupWithManager(mgr); err != nil {
				setupLog.Error(err, "unable to create controller", "controller", "TaskWorkload")
				os.Exit(1)
			}
		}

		if controllerConfig.IncludeStatefulsetRunner {
			logger := ctrl.Log.WithName("controllers").WithName("AppWorkload")
			if err = statefulsetcontrollers.NewAppWorkloadReconciler(
				mgr.GetClient(),
				mgr.GetScheme(),
				statefulsetcontrollers.NewAppWorkloadToStatefulsetConverter(mgr.GetScheme()),
				statefulsetcontrollers.NewPDBUpdater(mgr.GetClient()),
				logger,
			).SetupWithManager(mgr); err != nil {
				setupLog.Error(err, "unable to create controller", "controller", "AppWorkload")
				os.Exit(1)
			}

			logger = ctrl.Log.WithName("controllers").WithName("RunnerInfo")
			if err = statefulsetcontrollers.NewRunnerInfoReconciler(
				mgr.GetClient(),
				mgr.GetScheme(),
				logger,
			).SetupWithManager(mgr); err != nil {
				setupLog.Error(err, "unable to create controller", "controller", "RunnerInfo")
				os.Exit(1)
			}
		}

		if controllerConfig.IncludeContourRouter {
			if err = (networkingcontrollers.NewCFRouteReconciler(
				mgr.GetClient(),
				mgr.GetScheme(),
				ctrl.Log.WithName("controllers").WithName("CFRoute"),
				controllerConfig,
			)).SetupWithManager(mgr); err != nil {
				setupLog.Error(err, "unable to create controller", "controller", "CFRoute")
				os.Exit(1)
			}
		}

	}

	// Setup webhooks with manager

	if os.Getenv("ENABLE_WEBHOOKS") != "false" {
		if err = (&korifiv1alpha1.CFApp{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "CFApp")
			os.Exit(1)
		}

		(&workloads.AppRevWebhook{}).SetupWebhookWithManager(mgr)

		if err = (&korifiv1alpha1.CFPackage{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "CFPackage")
			os.Exit(1)
		}
		if err = (&korifiv1alpha1.CFBuild{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "CFBuild")
			os.Exit(1)
		}

		if err = korifiv1alpha1.NewCFProcessDefaulter(
			controllerConfig.CFProcessDefaults.MemoryMB,
			controllerConfig.CFProcessDefaults.DiskQuotaMB,
			*controllerConfig.CFProcessDefaults.Timeout,
		).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "CFProcess")
			os.Exit(1)
		}

		if err = workloads.NewCFAppValidator(
			webhooks.NewDuplicateValidator(coordination.NewNameRegistry(mgr.GetClient(), workloads.AppEntityType)),
		).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "CFApp")
			os.Exit(1)
		}

		if err = networking.NewCFRouteValidator(
			webhooks.NewDuplicateValidator(coordination.NewNameRegistry(mgr.GetClient(), networking.RouteEntityType)),
			controllerConfig.CFRootNamespace,
			mgr.GetClient(),
		).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "CFRoute")
			os.Exit(1)
		}

		if err = services.NewCFServiceInstanceValidator(
			webhooks.NewDuplicateValidator(coordination.NewNameRegistry(mgr.GetClient(), services.ServiceInstanceEntityType)),
		).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "CFServiceInstance")
			os.Exit(1)
		}

		if err = services.NewCFServiceBindingValidator(
			webhooks.NewDuplicateValidator(coordination.NewNameRegistry(mgr.GetClient(), services.ServiceBindingEntityType)),
		).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "CFServiceBinding")
			os.Exit(1)
		}

		if err = networking.NewCFDomainValidator(
			mgr.GetClient(),
		).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "CFDomain")
			os.Exit(1)
		}

		if err = workloads.NewCFOrgValidator(
			webhooks.NewDuplicateValidator(coordination.NewNameRegistry(mgr.GetClient(), workloads.CFOrgEntityType)),
			webhooks.NewPlacementValidator(mgr.GetClient(), controllerConfig.CFRootNamespace),
		).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "CFOrg")
			os.Exit(1)
		}

		if err = workloads.NewCFSpaceValidator(
			webhooks.NewDuplicateValidator(coordination.NewNameRegistry(mgr.GetClient(), workloads.CFSpaceEntityType)),
			webhooks.NewPlacementValidator(mgr.GetClient(), controllerConfig.CFRootNamespace),
		).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "CFSpace")
			os.Exit(1)
		}

		if err = (&korifiv1alpha1.CFRoute{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "CFRoute")
			os.Exit(1)
		}

		if err = workloads.NewCFTaskDefaulter(controllerConfig.CFProcessDefaults).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "CFTask")
			os.Exit(1)
		}

		if err = workloads.NewCFTaskValidator().SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "CFTask")
			os.Exit(1)
		}

		versionwebhook.NewVersionWebhook(version.Version).SetupWebhookWithManager(mgr)
		controllersfinalizer.NewControllersFinalizerWebhook().SetupWebhookWithManager(mgr)

		if err = workloads.NewCFPackageValidator().SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "CFPackage")
			os.Exit(1)
		}

		if controllerConfig.IncludeStatefulsetRunner {
			if err = statesetfulrunnerv1.NewSTSPodDefaulter().SetupWebhookWithManager(mgr); err != nil {
				setupLog.Error(err, "unable to create webhook", "webhook", "Pod")
				os.Exit(1)
			}
		}

		if controllerConfig.IncludeKpackImageBuilder {
			kpackimagebuilderfinalizer.NewKpackImageBuilderFinalizerWebhook().SetupWebhookWithManager(mgr)
		}

		if err = mgr.AddReadyzCheck("readyz", mgr.GetWebhookServer().StartedChecker()); err != nil {
			setupLog.Error(err, "unable to set up ready check")
			os.Exit(1)
		}
	} else {
		setupLog.Info("skipping webhook setup because ENABLE_WEBHOOKS set to false.")
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}

	eventChan := make(chan string)
	go func() {
		setupLog.Info("starting to watch config file at "+configPath+" for logger level changes", "currentLevel", atomicLevel.Level())
		if err2 := tools.WatchForConfigChangeEvents(context.Background(), configPath, setupLog, eventChan); err2 != nil {
			setupLog.Error(err2, "error watching logging config")
			os.Exit(1)
		}
	}()

	go tools.SyncLogLevel(context.Background(), setupLog, eventChan, atomicLevel, config.GetLogLevelFromPath)

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
