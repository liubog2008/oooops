package flow

type ControllerOptions struct {
}

// Controller defines controller to manage flow lifecycle and generate jobs
// It will do these things:
// - Watch flow and generate job of current stage
// - Go to next stage if job succeed
// - Mark as failed if job of current stage is failed
type Controller struct {
}
