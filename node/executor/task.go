package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type Task struct {
	ID            string
	JobID         string
	ShardID       string
	Status        string
	Progress      float64
	CheckpointCID string
	TaskType      string // "training", "inference", etc.
	ModelPath     string // IPFS CID or local path
	DatasetPath   string // IPFS CID or local path
	InputData     []byte // Input data for inference tasks
	OutputData    []byte // Output data for inference tasks
	CreatedAt     time.Time
	StartedAt     *time.Time
	CompletedAt   *time.Time
	Error         error
}

type Executor struct {
	resourceManager   interface{}
	tasks             map[string]*Task
	trainingExecutor  *TrainingExecutor
	inferenceExecutor *InferenceExecutor
	workDir           string
	ipfsAPIURL        string
	mu                sync.RWMutex
	ctx               context.Context
	cancel            context.CancelFunc
}

func NewExecutor(resourceManager interface{}) *Executor {
	ctx, cancel := context.WithCancel(context.Background())
	return &Executor{
		resourceManager: resourceManager,
		tasks:          make(map[string]*Task),
		workDir:        "/tmp/atlas-tasks",
		ipfsAPIURL:     "/ip4/127.0.0.1/tcp/5001",
		ctx:            ctx,
		cancel:         cancel,
	}
}

// SetWorkDir sets the working directory for tasks
func (e *Executor) SetWorkDir(workDir string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.workDir = workDir
}

// SetIPFSAPIURL sets the IPFS API URL
func (e *Executor) SetIPFSAPIURL(ipfsAPIURL string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.ipfsAPIURL = ipfsAPIURL
}

// InitializeTrainingExecutor initializes the training executor
func (e *Executor) InitializeTrainingExecutor() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.trainingExecutor == nil {
		e.trainingExecutor = NewTrainingExecutor(e, e.workDir, e.ipfsAPIURL)
	}
}

// InitializeInferenceExecutor initializes the inference executor
func (e *Executor) InitializeInferenceExecutor() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.inferenceExecutor == nil {
		e.inferenceExecutor = NewInferenceExecutor(e, e.workDir, e.ipfsAPIURL)
	}
}

func (e *Executor) Start(ctx context.Context) error {
	// Initialize training executor if not already initialized
	e.InitializeTrainingExecutor()
	
	// Use provided context or executor's context
	execCtx := ctx
	if execCtx == nil {
		execCtx = e.ctx
	}
	
	ticker := time.NewTicker(5 * time.Second) // Check every 5 seconds
	defer ticker.Stop()

	for {
		select {
		case <-execCtx.Done():
			return execCtx.Err()
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			e.processTasks(execCtx)
		}
	}
}

// AddTask adds a new task to the executor
func (e *Executor) AddTask(task *Task) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	if task.ID == "" {
		return fmt.Errorf("task ID cannot be empty")
	}
	
	if _, exists := e.tasks[task.ID]; exists {
		return fmt.Errorf("task %s already exists", task.ID)
	}
	
	if task.Status == "" {
		task.Status = "pending"
	}
	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now()
	}
	
	e.tasks[task.ID] = task
	return nil
}

// GetTask retrieves a task by ID
func (e *Executor) GetTask(taskID string) (*Task, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	
	task, ok := e.tasks[taskID]
	if !ok {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}
	
	return task, nil
}

func (e *Executor) processTasks(ctx context.Context) {
	e.mu.RLock()
	tasks := make([]*Task, 0, len(e.tasks))
	for _, task := range e.tasks {
		tasks = append(tasks, task)
	}
	e.mu.RUnlock()
	
	for _, task := range tasks {
		e.mu.Lock()
		status := task.Status
		e.mu.Unlock()
		
		switch status {
		case "pending":
			// Start pending tasks
			e.startTask(ctx, task)
		case "in_progress":
			// Tasks in progress are handled by goroutines
			// This is just a check to see if they're still running
			continue
		case "completed", "failed", "paused":
			// Skip completed, failed, or paused tasks
			continue
		}
	}
}

func (e *Executor) startTask(ctx context.Context, task *Task) {
	e.mu.Lock()
	if task.Status != "pending" {
		e.mu.Unlock()
		return
	}
	task.Status = "in_progress"
	now := time.Now()
	task.StartedAt = &now
	e.mu.Unlock()
	
	// Execute task in goroutine
	go e.executeTask(ctx, task)
}

func (e *Executor) executeTask(ctx context.Context, task *Task) {
	defer func() {
		e.mu.Lock()
		if task.Status == "in_progress" {
			if task.Error != nil {
				task.Status = "failed"
			} else {
				task.Status = "completed"
				now := time.Now()
				task.CompletedAt = &now
				task.Progress = 1.0
			}
		}
		e.mu.Unlock()
	}()
	
	// Execute based on task type
	switch task.TaskType {
	case "training":
		e.InitializeTrainingExecutor()
		e.executeTrainingTask(ctx, task)
	case "inference":
		e.InitializeInferenceExecutor()
		e.executeInferenceTask(ctx, task)
	default:
		// Default to training if not specified
		if task.TaskType == "" {
			task.TaskType = "training"
			e.executeTrainingTask(ctx, task)
		} else {
			task.Error = fmt.Errorf("unknown task type: %s", task.TaskType)
			return
		}
	}
}

func (e *Executor) executeTrainingTask(ctx context.Context, task *Task) {
	e.mu.RLock()
	trainingExecutor := e.trainingExecutor
	e.mu.RUnlock()
	
	if trainingExecutor == nil {
		task.Error = fmt.Errorf("training executor not initialized")
		return
	}
	
	// Update progress during execution
	e.mu.Lock()
	task.Progress = 0.1 // Started
	e.mu.Unlock()
	
	// Execute training
	if err := trainingExecutor.ExecuteTraining(ctx, task, task.ModelPath, task.DatasetPath); err != nil {
		task.Error = fmt.Errorf("training execution failed: %w", err)
		return
	}
	
	e.mu.Lock()
	task.Progress = 1.0 // Completed
	e.mu.Unlock()
}

func (e *Executor) executeInferenceTask(ctx context.Context, task *Task) {
	e.mu.RLock()
	inferenceExecutor := e.inferenceExecutor
	e.mu.RUnlock()
	
	if inferenceExecutor == nil {
		task.Error = fmt.Errorf("inference executor not initialized")
		return
	}
	
	if len(task.InputData) == 0 {
		task.Error = fmt.Errorf("input data is required for inference tasks")
		return
	}
	
	e.mu.Lock()
	task.Progress = 0.1
	e.mu.Unlock()
	
	output, err := inferenceExecutor.ExecuteInference(ctx, task, task.ModelPath, task.InputData)
	if err != nil {
		task.Error = fmt.Errorf("inference execution failed: %w", err)
		return
	}
	
	outputJSON, err := json.Marshal(output)
	if err != nil {
		task.Error = fmt.Errorf("failed to marshal output: %w", err)
		return
	}
	
	e.mu.Lock()
	task.OutputData = outputJSON
	task.Progress = 1.0
	e.mu.Unlock()
}

// ExecuteTask manually starts a task execution
func (e *Executor) ExecuteTask(taskID string) error {
	task, err := e.GetTask(taskID)
	if err != nil {
		return err
	}
	
	e.mu.Lock()
	if task.Status != "pending" && task.Status != "paused" {
		e.mu.Unlock()
		return fmt.Errorf("task %s is not in pending or paused state (current: %s)", taskID, task.Status)
	}
	task.Status = "in_progress"
	now := time.Now()
	task.StartedAt = &now
	e.mu.Unlock()
	
	// Execute in background
	go e.executeTask(e.ctx, task)
	return nil
}

// StopTask stops/pauses a running task
func (e *Executor) StopTask(taskID string) error {
	task, err := e.GetTask(taskID)
	if err != nil {
		return err
	}
	
	e.mu.Lock()
	defer e.mu.Unlock()
	
	if task.Status == "in_progress" {
		task.Status = "paused"
		return nil
	}
	
	return fmt.Errorf("task %s is not in progress (current: %s)", taskID, task.Status)
}

// ListTasks returns all tasks
func (e *Executor) ListTasks() []*Task {
	e.mu.RLock()
	defer e.mu.RUnlock()
	
	tasks := make([]*Task, 0, len(e.tasks))
	for _, task := range e.tasks {
		tasks = append(tasks, task)
	}
	
	return tasks
}

// Stop stops the executor and cancels all running tasks
func (e *Executor) Stop() {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	if e.cancel != nil {
		e.cancel()
	}
	
	// Mark all in-progress tasks as paused
	for _, task := range e.tasks {
		if task.Status == "in_progress" {
			task.Status = "paused"
		}
	}
}

