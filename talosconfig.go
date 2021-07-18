package taloscdk

type TalosRole string

const (
	CONTROL_PLANE TalosRole = "control_plane"
	WORKER        TalosRole = "worker"
)
