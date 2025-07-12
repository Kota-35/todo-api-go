package task

// type CreateTodoRequest struct {
// 	Title       string `json:"title" binding:"required"`
// 	Description string `json:"description" binding:"required"`
// 	Status      string `json:"status" binding:"required"`
// }

// func (h *Handler) Create(c *gin.Context) {
// 	var req CreateTodoRequest

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	todo, err := h.taskService.Create(req.Title, req.Description, req.Status)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, todo)
// }
