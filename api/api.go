package api

import (
	"github.com/labstack/echo/v4"
	"miami-back-end/project"
	"miami-back-end/members"
	"miami-back-end/meeting"
	"miami-back-end/stage"
	"miami-back-end/pilot-questionnaire"
	"miami-back-end/operation-training"
	"miami-back-end/questionnaire-sign-off"
	"miami-back-end/revised-questionnaire"
	"miami-back-end/data-collection"
	"miami-back-end/qaqc-training"
	"miami-back-end/qaqc-collection"
	"miami-back-end/de-training"
	"miami-back-end/de-collection"
	"miami-back-end/da-collection"
	"miami-back-end/report-submission-project-close"
	"miami-back-end/project-review"

)

func InitRoutes(e *echo.Echo) {
	apiGroup := e.Group("/miami-api")
	apiGroup.Use(Auth())
	
	// สร้าง sub-group /stage
	stageGroup := apiGroup.Group("/stage")
	stageService := stage.NewStageService()
	stageRoute := NewStageRoute(stageService)
	stageRoute.Group(stageGroup)

	// สร้าง sub-group /member
	memberGroup := apiGroup.Group("/member")
	memberService := members.NewMemberService()
	memberRoute := NewMemberRoute(memberService,stageService)
	memberRoute.Group(memberGroup)

	// สร้าง sub-group /project
	projectGroup := apiGroup.Group("/project")
	projectService := project.NewService()
	projectRoute := NewProjectRoute(projectService,memberService)
	projectRoute.Group(projectGroup)
	
	// สร้าง sub-group /meeting
	meetingGroup := apiGroup.Group("/meeting")
	meetingService := meeting.NewMeetingService()
	meetingRoute := NewMeetingRoute(meetingService, memberService, projectService, stageService)
	meetingRoute.Group(meetingGroup)

	// สร้าง sub-group /pilot-questionnaire
	pilotQuestionnaireGroup := apiGroup.Group("/pilot-questionnaire")
	pilotQuestionnaireService := pilotquestionnaire.NewPilotQuestionnaireService()
	pilotQuestionnaireRoute := NewPilotQuesttionaireRoute(pilotQuestionnaireService,projectService,memberService,stageService)
	pilotQuestionnaireRoute.Group(pilotQuestionnaireGroup)

	// สร้าง sub-group /questionnaire-sign-off
	questionnaireSignOffGroup := apiGroup.Group("/questionnaire-sign-off")
	questionnaireSignOffService := questionnairesignoff.NewQuestionnaireSignOffService()
	questionnaireSignOffRoute := NewQuestionnaireSORoute(questionnaireSignOffService,memberService,projectService,stageService,pilotQuestionnaireService)
	questionnaireSignOffRoute.Group(questionnaireSignOffGroup)

	// สร้าง sub-group /operation-training
	trainingGroup := apiGroup.Group("/operation-training")
	trainingService := operationtraining.NewOperationTrainingService()
	trainingRoute := NewTrainingRoute(trainingService,stageService,memberService,projectService,pilotQuestionnaireService)
	trainingRoute.Group(trainingGroup)

	// สร้าง sub-group /revised-questionnaire
	revisedQuestionnaireGroup := apiGroup.Group("/revised-questionnaire")
	revisedQuestionnaireService := revisedquestionnaire.NewRQService()
	revisedQuestionnaireRoute := NewRQRoute(revisedQuestionnaireService,pilotQuestionnaireService,stageService,projectService,memberService)
	revisedQuestionnaireRoute.Group(revisedQuestionnaireGroup)

	//data-collection
	dataCollectionGroup := apiGroup.Group("/data-collection")
	dataCollectionService := datacollection.NewDataCollectionService()
	dataCollectionRoute := NewDataCollectionRoute(dataCollectionService,projectService,memberService,stageService)
	dataCollectionRoute.Group(dataCollectionGroup)

	//qaqc training
	qaqcTrainingGroup := apiGroup.Group("/qaqc-training")
	qaqcTrainingService := qaqctraining.NewQAQCService()
	qaqcTrainingRoute := NewQAQCTrainingRoute(qaqcTrainingService,memberService,projectService,stageService,pilotQuestionnaireService)
	qaqcTrainingRoute.Group(qaqcTrainingGroup)

	//qaqc collection
	qaqcCollectionGroup := apiGroup.Group("/qaqc-collection")
	qaqcCollectionService := qaqccollection.NewQaqcCollectionService()
	qaqcCollectionRoute := NewQAqcRoute(qaqcCollectionService,projectService,memberService,stageService)
	qaqcCollectionRoute.Group(qaqcCollectionGroup)

	// de training
	deTrainingGroup := apiGroup.Group("/de-training")
	deTrainingService := detraining.NewDeTrainingService()
	deTrainingRoute := NewDeTrainingRoute(deTrainingService,memberService,projectService,stageService,pilotQuestionnaireService)
	deTrainingRoute.Group(deTrainingGroup)

	// de collection
	deCollectionGroup := apiGroup.Group("/de-collection")
	deCollectionService := decollection.NewDeCollectionService()
	deCollectionRoute := NewDeCRoute(deCollectionService,projectService,memberService,stageService)
	deCollectionRoute.Group(deCollectionGroup)

	// da collection
	daCollectionGroup := apiGroup.Group("/da-collection")
	daCollectionService := dacollection.NewDaCollectionService()
	daCollectionRoute := NewDaCRoute(daCollectionService,projectService,memberService,stageService)
	daCollectionRoute.Group(daCollectionGroup)

	// report submission and project close
	reportSubmissionGroup := apiGroup.Group("/report-submission-project-close")
	reportSubmissionService := reportsubmissionprojectclose.NewReportSubmissionService()
	reportSubmissionRoute := NewRsPcRoute(reportSubmissionService,projectService,memberService,stageService)
	reportSubmissionRoute.Group(reportSubmissionGroup)

	// Project Review
	projectReviewGroup := apiGroup.Group("/project-review")
	projectReviewService := projectreview.NewProjectReviewService()
	projectReviewRoute := NewProjectReviewRoute(projectReviewService,memberService,projectService,stageService)
	projectReviewRoute.Group(projectReviewGroup)


	









	

}