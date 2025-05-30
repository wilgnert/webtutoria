-- name: ResetSSC :exec
delete from StudentSubjectCompletion;
-- name: ResetSD :exec
delete from StudentDiscords;
-- name: ResetTD :exec
delete from TutorDiscords;
-- name: ResetST :exec
delete from StudentTutor;
-- name: ResetSC :exec
delete from SubjectCategory;
-- name: ResetCategories :exec
delete from Categories;
-- name: ResetStudents :exec
delete from Students;
-- name: ResetSubjects :exec
delete from Subjects;
-- name: ResetTutors :exec
delete from Tutors;