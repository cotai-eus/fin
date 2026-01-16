package support

import "errors"

var (
	// ErrTicketNotFound is returned when a support ticket is not found
	ErrTicketNotFound = errors.New("ticket not found")

	// ErrMessageNotFound is returned when a ticket message is not found
	ErrMessageNotFound = errors.New("message not found")

	// ErrUnauthorizedAccess is returned when user tries to access another user's ticket
	ErrUnauthorizedAccess = errors.New("unauthorized access to ticket")

	// ErrInvalidCategory is returned when ticket category is invalid
	ErrInvalidCategory = errors.New("invalid ticket category")

	// ErrInvalidPriority is returned when ticket priority is invalid
	ErrInvalidPriority = errors.New("invalid ticket priority")

	// ErrInvalidStatus is returned when ticket status is invalid
	ErrInvalidStatus = errors.New("invalid ticket status")

	// ErrInvalidStatusTransition is returned when status transition is not allowed
	ErrInvalidStatusTransition = errors.New("invalid status transition")

	// ErrEmptySubject is returned when ticket subject is empty
	ErrEmptySubject = errors.New("ticket subject cannot be empty")

	// ErrEmptyDescription is returned when ticket description is empty
	ErrEmptyDescription = errors.New("ticket description cannot be empty")

	// ErrEmptyMessage is returned when ticket message is empty
	ErrEmptyMessage = errors.New("message cannot be empty")

	// ErrTicketClosed is returned when trying to add message to closed ticket
	ErrTicketClosed = errors.New("cannot add message to closed ticket")

	// ErrSubjectTooLong is returned when subject exceeds maximum length
	ErrSubjectTooLong = errors.New("subject exceeds maximum length of 255 characters")

	// ErrDescriptionTooShort is returned when description is too short
	ErrDescriptionTooShort = errors.New("description must be at least 10 characters")

	// ErrMessageTooShort is returned when message is too short
	ErrMessageTooShort = errors.New("message must be at least 5 characters")
)
