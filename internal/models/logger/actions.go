package logger

type UserAction string
type UserPunishment string
type CommonAction string

// ! User punishments
const (
	Ban   UserPunishment = "has banned"
	Unban UserPunishment = "has unbanned"
)

// ! Common actions
const (
	CreateRank CommonAction = "has created rank"

	SearchByUsername CommonAction = "searched by username"
	SearchByUserID   CommonAction = "searched by user ID"
	SearchByAllUsers CommonAction = "searched all users"
	SearchLogs       CommonAction = "searched logs"

	SetStaffRank     UserAction = "has set admin perm"
	SetDeveloperRank UserAction = "has set developer perm"
	RestoreUser      UserAction = "has restored"
	Create           UserAction = "has created"
	ChangeFlags      UserAction = "has changed flags"

	SoftDelete           UserPunishment = "has soft-deleted"
	HardDelete           UserPunishment = "has hard-deleted"
	TriedToDeleteManager UserPunishment = "has tried to delete manager's account and action has stopped"
)
