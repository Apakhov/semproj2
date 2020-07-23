package main

import (
	"net/http"
	"semprojdb/db/psql"
	"semprojdb/handler"
	"semprojdb/logger"

	"github.com/jackc/pgx/v4"

	"github.com/gorilla/mux"
)

func main() {
	lg := logger.NewLogger("", "", "log.log", pgx.LogLevelDebug)
	db := psql.NewDB(&psql.Config{
		Host:     "localhost",
		Database: "study_ev_table",
		User:     "postgres",
		Password: "docker",
	}, lg.SubLog("PSQL"))
	err := db.Connect()
	if err != nil {
		lg.Fatal(err)
	}
	h := handler.NewHandler(db, lg.SubLog("HANDLER"))

	router := mux.NewRouter()

	router.HandleFunc("/faculty", h.HandleFacultyCreate).Methods(http.MethodPost)
	router.HandleFunc("/faculty", h.HandleFacultyRead).Methods(http.MethodGet)
	router.HandleFunc("/faculty", h.HandleFacultyUpdate).Methods(http.MethodPut)
	router.HandleFunc("/faculty", h.HandleFacultyDelete).Methods(http.MethodDelete)

	router.HandleFunc("/cathedra", h.HandleCathedraCreate).Methods(http.MethodPost)
	router.HandleFunc("/cathedra", h.HandleCathedraRead).Methods(http.MethodGet)
	router.HandleFunc("/cathedra", h.HandleCathedraUpdate).Methods(http.MethodPut)
	router.HandleFunc("/cathedra", h.HandleCathedraDelete).Methods(http.MethodDelete)

	router.HandleFunc("/subject", h.HandleSubjectCreate).Methods(http.MethodPost)
	router.HandleFunc("/subject", h.HandleSubjectRead).Methods(http.MethodGet)
	router.HandleFunc("/subject", h.HandleSubjectUpdate).Methods(http.MethodPut)
	router.HandleFunc("/subject", h.HandleSubjectDelete).Methods(http.MethodDelete)

	router.HandleFunc("/teacher", h.HandleTeacherCreate).Methods(http.MethodPost)
	router.HandleFunc("/teacher", h.HandleTeacherRead).Methods(http.MethodGet)
	router.HandleFunc("/teacher", h.HandleTeacherUpdate).Methods(http.MethodPut)
	router.HandleFunc("/teacher", h.HandleTeacherDelete).Methods(http.MethodDelete)

	router.HandleFunc("/st_group", h.HandleStGroupCreate).Methods(http.MethodPost)
	router.HandleFunc("/st_group", h.HandleStGroupRead).Methods(http.MethodGet)
	router.HandleFunc("/st_group", h.HandleStGroupUpdate).Methods(http.MethodPut)
	router.HandleFunc("/st_group", h.HandleStGroupDelete).Methods(http.MethodDelete)

	router.HandleFunc("/student", h.HandleStudentCreate).Methods(http.MethodPost)
	router.HandleFunc("/student", h.HandleStudentRead).Methods(http.MethodGet)
	router.HandleFunc("/student", h.HandleStudentUpdate).Methods(http.MethodPut)
	router.HandleFunc("/student", h.HandleStudentDelete).Methods(http.MethodDelete)

	router.HandleFunc("/stgr", h.HandlePutStudentToGroup).Methods(http.MethodPost)
	router.HandleFunc("/stgr/st", h.HandleGetStudentsFromGroup).Methods(http.MethodGet)
	router.HandleFunc("/stgr/gr", h.HandleGetGroupsOfStudent).Methods(http.MethodGet)

	router.HandleFunc("/course", h.HandleCourseCreate).Methods(http.MethodPost)
	router.HandleFunc("/course", h.HandleCourseRead).Methods(http.MethodGet)
	router.HandleFunc("/course", h.HandleCourseUpdate).Methods(http.MethodPut)
	router.HandleFunc("/course", h.HandleCourseDelete).Methods(http.MethodDelete)

	router.HandleFunc("/teachercourse", h.HandleAssignTeacherToCourse).Methods(http.MethodPost)
	router.HandleFunc("/teachercourse", h.HandleUnassignTeacherFromCourse).Methods(http.MethodDelete)
	router.HandleFunc("/course/teachers", h.HandleGetTeachersFromCourse).Methods(http.MethodGet)
	router.HandleFunc("/teacher/courses", h.HandleGetCoursesOfTeacher).Methods(http.MethodGet)

	router.HandleFunc("/attendance", h.HandleAttendanceCreate).Methods(http.MethodPost)
	router.HandleFunc("/attendance", h.HandleAttendanceRead).Methods(http.MethodGet)
	router.HandleFunc("/attendance", h.HandleAttendanceUpdate).Methods(http.MethodPut)
	router.HandleFunc("/attendance", h.HandleAttendanceDelete).Methods(http.MethodDelete)

	router.HandleFunc("/mark", h.HandleMarkCreate).Methods(http.MethodPost)
	router.HandleFunc("/mark", h.HandleMarkRead).Methods(http.MethodGet)
	router.HandleFunc("/mark", h.HandleMarkUpdate).Methods(http.MethodPut)
	router.HandleFunc("/mark", h.HandleMarkDelete).Methods(http.MethodDelete)

	router.HandleFunc("/exam", h.HandleExamCreate).Methods(http.MethodPost)
	router.HandleFunc("/exam", h.HandleExamRead).Methods(http.MethodGet)
	router.HandleFunc("/exam", h.HandleExamUpdate).Methods(http.MethodPut)
	router.HandleFunc("/exam", h.HandleExamDelete).Methods(http.MethodDelete)

	http.Handle("/", router)
	lg.Infof("Server is listening on 8181...")
	http.ListenAndServe(":8181", nil)
}
