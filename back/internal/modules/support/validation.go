package support

import "strings"

// ValidateCategory validates ticket category
func ValidateCategory(category string) error {
	if category == "" {
		return ErrInvalidCategory
	}

	category = strings.ToLower(strings.TrimSpace(category))
	if !ValidCategories[category] {
		return ErrInvalidCategory
	}

	return nil
}

// ValidatePriority validates ticket priority
func ValidatePriority(priority string) error {
	if priority == "" {
		return ErrInvalidPriority
	}

	priority = strings.ToLower(strings.TrimSpace(priority))
	if !ValidPriorities[priority] {
		return ErrInvalidPriority
	}

	return nil
}

// ValidateStatus validates ticket status
func ValidateStatus(status string) error {
	if status == "" {
		return ErrInvalidStatus
	}

	status = strings.ToLower(strings.TrimSpace(status))
	if !ValidStatuses[status] {
		return ErrInvalidStatus
	}

	return nil
}

// ValidateStatusTransition validates if a status transition is allowed
func ValidateStatusTransition(currentStatus, newStatus string) error {
	if currentStatus == newStatus {
		return nil // No transition needed
	}

	allowedTransitions, exists := AllowedStatusTransitions[currentStatus]
	if !exists {
		return ErrInvalidStatusTransition
	}

	for _, allowed := range allowedTransitions {
		if allowed == newStatus {
			return nil
		}
	}

	return ErrInvalidStatusTransition
}

// ValidateSubject validates ticket subject
func ValidateSubject(subject string) error {
	subject = strings.TrimSpace(subject)

	if subject == "" {
		return ErrEmptySubject
	}

	if len(subject) > 255 {
		return ErrSubjectTooLong
	}

	return nil
}

// ValidateDescription validates ticket description
func ValidateDescription(description string) error {
	description = strings.TrimSpace(description)

	if description == "" {
		return ErrEmptyDescription
	}

	if len(description) < 10 {
		return ErrDescriptionTooShort
	}

	return nil
}

// ValidateMessage validates ticket message
func ValidateMessage(message string) error {
	message = strings.TrimSpace(message)

	if message == "" {
		return ErrEmptyMessage
	}

	if len(message) < 5 {
		return ErrMessageTooShort
	}

	return nil
}

// ValidateCreateTicketRequest validates the create ticket request
func ValidateCreateTicketRequest(req CreateTicketRequest) error {
	if err := ValidateCategory(req.Category); err != nil {
		return err
	}

	if err := ValidatePriority(req.Priority); err != nil {
		return err
	}

	if err := ValidateSubject(req.Subject); err != nil {
		return err
	}

	if err := ValidateDescription(req.Description); err != nil {
		return err
	}

	return nil
}
