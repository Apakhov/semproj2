package testing

import (
	"net/http"
	"semprojdb/ftests/test"
	"semprojdb/handler"
	"testing"
	"time"
)

func testCRUD(
	t *testing.T,
	path string,
	new func() interface{},
	update func(interface{}) interface{},
) (interface{}, interface{}) {
	v := new()
	tc := test.TestCase{
		Path:      "/faculty",
		Method:    http.MethodPost,
		GetParams: nil,
		Body:      v,
		Response: test.WithError{
			Err: "",
			Val: v,
		},
	}
	v2 := new()
	tc = test.TestCase{
		Path:      "/faculty",
		Method:    http.MethodPost,
		GetParams: nil,
		Body:      v2,
		Response: test.WithError{
			Err: "",
			Val: v2,
		},
	}

	vUPD := new()
	tc.Run(t, &vUPD)

	vUPD = update(vUPD)
	tc = test.TestCase{
		Path:      "/faculty",
		Method:    http.MethodPut,
		GetParams: nil,
		Body:      vUPD,
		Response: test.WithError{
			Err: "",
			Val: vUPD,
		},
	}

	tc.Run(t, &vUPD)
	return v, v2
}
func NewFaculty(t *testing.T) handler.Faculty {
	faculty := handler.Faculty{
		ShortName: test.UStr(),
		FullName:  test.UStr(),
	}
	tc := test.TestCase{
		Path:      "/faculty",
		Method:    http.MethodPost,
		GetParams: nil,
		Body:      faculty,
		Response: test.WithError{
			Err: "",
			Val: faculty,
		},
	}

	tc.Run(t, &faculty)
	return faculty
}
func UpdateAny(t *testing.T, v interface{}, path string) {
	tc := test.TestCase{
		Path:      path,
		Method:    http.MethodPut,
		GetParams: nil,
		Body:      v,
		Response: test.WithError{
			Err: "",
			Val: v,
		},
	}

	tc.Run(t, &v)
}
func DelFaculty(t *testing.T, v *handler.Faculty) {

}
func NewCathedra(t *testing.T, facultyID int64) handler.Cathedra {
	cathedra := handler.Cathedra{
		ShortName: test.UStr(),
		FullName:  test.UStr(),
		FacultyID: facultyID,
	}
	tc := test.TestCase{
		Path:      "/cathedra",
		Method:    http.MethodPost,
		GetParams: nil,
		Body:      cathedra,
		Response: test.WithError{
			Err: "",
			Val: cathedra,
		},
	}

	tc.Run(t, &cathedra)
	return cathedra
}
func DelCathedra(t *testing.T, v *handler.Cathedra) {

}
func NewSubject(t *testing.T, cathedraID int64) handler.Subject {
	subject := handler.Subject{
		ShortName:  test.UStr(),
		FullName:   test.UStr(),
		CathedraID: cathedraID,
	}
	tc := test.TestCase{
		Path:      "/subject",
		Method:    http.MethodPost,
		GetParams: nil,
		Body:      subject,
		Response: test.WithError{
			Err: "",
			Val: subject,
		},
	}

	tc.Run(t, &subject)
	return subject
}
func DelSubject(t *testing.T, v *handler.Subject) {

}
func NewTeacher(t *testing.T, cathedraID int64) handler.Teacher {
	v := handler.Teacher{
		ContractID: test.UStr(),
		FirsName:   test.UStr(),
		LastName:   test.UStr(),
		Email:      test.UEmail(),
		CathedraID: cathedraID,
	}
	tc := test.TestCase{
		Path:      "/teacher",
		Method:    http.MethodPost,
		GetParams: nil,
		Body:      v,
		Response: test.WithError{
			Err: "",
			Val: v,
		},
	}

	tc.Run(t, &v)
	return v
}
func DelTeacher(t *testing.T, v *handler.Teacher) {

}
func NewStGroup(t *testing.T, teacherID int64) handler.StGroup {
	v := handler.StGroup{
		GroupID:   test.UStr(),
		BeginD:    time.Now(),
		EndD:      time.Now().Add(180 * 24 * time.Hour),
		TeacherID: teacherID,
	}
	tc := test.TestCase{
		Path:      "/st_group",
		Method:    http.MethodPost,
		GetParams: nil,
		Body:      v,
		Response: test.WithError{
			Err: "",
			Val: v,
		},
	}

	tc.Run(t, &v)
	return v
}
func DelStGroup(t *testing.T, v *handler.StGroup) {

}
func NewStudent(t *testing.T) handler.Student {
	v := handler.Student{
		StudyID:  test.UStr(),
		FirsName: test.UStr(),
		LastName: test.UStr(),
		Email:    test.UEmail(),
	}
	tc := test.TestCase{
		Path:      "/student",
		Method:    http.MethodPost,
		GetParams: nil,
		Body:      v,
		Response: test.WithError{
			Err: "",
			Val: v,
		},
	}

	tc.Run(t, &v)
	return v
}
func DelStudent(t *testing.T, v *handler.Student) {

}
func NewCourse(t *testing.T, subjectID, stgroupID int64) handler.Course {
	v := handler.Course{
		ShortName: test.UStr(),
		FullName:  test.UStr(),
		Semester:  5,
		BeginD:    time.Now(),
		EndD:      time.Now().Add(180 * 24 * time.Hour),
		SubjectID: subjectID,
		StGroupID: stgroupID,
		Active:    true,
	}
	tc := test.TestCase{
		Path:      "/course",
		Method:    http.MethodPost,
		GetParams: nil,
		Body:      v,
		Response: test.WithError{
			Err: "",
			Val: v,
		},
	}

	tc.Run(t, &v)
	return v
}
func DelCourse(t *testing.T, v *handler.Course) {

}
func NewMark(t *testing.T, StudentID, TeacherID, CourseID int64) handler.Mark {
	v := handler.Mark{
		Date:      time.Now(),
		Points:    2,
		StudentID: StudentID,
		TeacherID: TeacherID,
		CourseID:  CourseID,
	}
	tc := test.TestCase{
		Path:      "/mark",
		Method:    http.MethodPost,
		GetParams: nil,
		Body:      v,
		Response: test.WithError{
			Err: "",
			Val: v,
		},
	}

	tc.Run(t, &v)
	return v
}
func NewAttendance(t *testing.T, StudentID, TeacherID, CourseID int64) handler.Attendance {
	v := handler.Attendance{
		Date:      time.Now(),
		StudentID: StudentID,
		TeacherID: TeacherID,
		CourseID:  CourseID,
	}
	tc := test.TestCase{
		Path:      "/attendance",
		Method:    http.MethodPost,
		GetParams: nil,
		Body:      v,
		Response: test.WithError{
			Err: "",
			Val: v,
		},
	}

	tc.Run(t, &v)
	return v
}
func NewExam(t *testing.T, StudentID, TeacherID, CourseID int64) handler.Exam {
	v := handler.Exam{
		Date:      time.Now(),
		Type:      "p",
		Points:    2,
		StudentID: StudentID,
		TeacherID: TeacherID,
		CourseID:  CourseID,
	}
	tc := test.TestCase{
		Path:      "/exam",
		Method:    http.MethodPost,
		GetParams: nil,
		Body:      v,
		Response: test.WithError{
			Err: "",
			Val: v,
		},
	}

	tc.Run(t, &v)
	return v
}

func TestAAA(t *testing.T) {
	// FACULTY -------------------------------------------------------
	faculty := NewFaculty(t)
	faculty.FullName = test.UStr()
	UpdateAny(t, &faculty, "/faculty")
	// FACULTY ^------------------------------------------------------

	// CATHEDRA ------------------------------------------------------
	cathedra := NewCathedra(t, faculty.ID)
	cathedra.FullName = test.UStr()
	UpdateAny(t, &cathedra, "/cathedra")
	// CATHEDRA ^-----------------------------------------------------

	// SUBJECT -------------------------------------------------------
	subject := NewSubject(t, cathedra.ID)
	subject.FullName = test.UStr()
	UpdateAny(t, &subject, "/subject")
	// SUBJECT ^------------------------------------------------------

	// TEACHER -------------------------------------------------------
	teacher := NewTeacher(t, cathedra.ID)
	teacher.FirsName = test.UStr()
	UpdateAny(t, &teacher, "/teacher")
	// TEACHER ^------------------------------------------------------

	// STGROUP -------------------------------------------------------
	stGroup := NewStGroup(t, teacher.ID)
	stGroup.GroupID = test.UStr()
	UpdateAny(t, &stGroup, "/st_group")
	// STGROUP ^------------------------------------------------------

	// STUDENT -------------------------------------------------------
	student := NewStudent(t)
	student.FirsName = test.UStr()
	UpdateAny(t, &student, "/student")
	// STUDENT ^------------------------------------------------------

	// COURSE --------------------------------------------------------
	course := NewCourse(t, subject.ID, stGroup.ID)
	course.FullName = test.UStr()
	UpdateAny(t, &course, "/course")
	// COURSE ^-------------------------------------------------------

	// MARK ----------------------------------------------------------
	mark := NewMark(t, student.ID, teacher.ID, course.ID)
	mark.Points = 3
	UpdateAny(t, &mark, "/mark")
	// MARK ^---------------------------------------------------------

	// ATTENDANCE ----------------------------------------------------
	attendance := NewAttendance(t, student.ID, teacher.ID, course.ID)
	attendance.Date = time.Now()
	UpdateAny(t, &attendance, "/attendance")
	// ATTENDANCE ^---------------------------------------------------

	// EXAM ----------------------------------------------------------
	exam := NewExam(t, student.ID, teacher.ID, course.ID)
	exam.Type = "np"
	exam.Points = 1
	UpdateAny(t, &exam, "/exam")
	// EXAM ^---------------------------------------------------------
}

// return handler.Faculty{
// 	ShortName: test.UStr(),
// 	FullName:  test.UStr(),
// }
// tc = test.TestCase{
// 	Path:      "/faculty",
// 	Method:    http.MethodPost,
// 	GetParams: nil,
// 	Body:      v,
// 	Response: test.WithError{
// 		Err: "",
// 		Val: v,
// 	},
// }

// nreturn handler.Faculty{}
// tc.Run(t, nv)

// nv.FullName = test.UStr()
// tc = test.TestCase{
// 	Path:      "/faculty",
// 	Method:    http.MethodPut,
// 	GetParams: nil,
// 	Body:      nv,
// 	Response: test.WithError{
// 		Err: "",
// 		Val: nv,
// 	},
// }
