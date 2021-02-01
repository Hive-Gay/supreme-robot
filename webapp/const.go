package webapp

type contextKey int

const SessionKey contextKey = 0
const UserKey contextKey = 1

const groupUWUCrew = "/UWU Crew"
const groupMailAdmin = "/Mail Admin"

var adminGroups = []string{
	groupMailAdmin,
}