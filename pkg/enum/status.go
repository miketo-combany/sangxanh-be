package enum

type Status string

const (
	Active   Status = "active"
	Inactive Status = "inactive"
)

func ToStatus(status bool) Status {
	if status {
		return Active
	}
	return Inactive
}
