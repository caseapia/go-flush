package logger

type LoggerAction string

const (
	// ! Admin actions
	Ban                  LoggerAction = "has banned"
	Unban                LoggerAction = "has unbanned"
	Create               LoggerAction = "has created"
	SoftDelete           LoggerAction = "has soft-deleted"
	HardDelete           LoggerAction = "has hard-deleted"
	RestoreUser          LoggerAction = "has restored"
	SetStaffRank         LoggerAction = "has set admin perm"
	SetDeveloperRank     LoggerAction = "has set developer perm"
	TriedToDeleteManager LoggerAction = "has tried to delete manager's account and action has stopped"
	CreateRank           LoggerAction = "has created rank"
	ChangeFlags          LoggerAction = "has changed flags"

	// ! Searches
	SearchByUsername LoggerAction = "searched by username"
	SearchByUserID   LoggerAction = "searched by user ID"
	SearchByAllUsers LoggerAction = "searched all users"
	SearchLogs       LoggerAction = "searched logs"
)
