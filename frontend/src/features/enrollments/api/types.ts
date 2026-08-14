export type EnrollmentStatus = 'active' | 'completed'

// Enrollment is serialized directly from the Go model, which has no JSON tags.
export interface EnrollmentDto {
  ID: string
  UserID: string
  CourseID: string
  Status: EnrollmentStatus
  StartedAt: string
  CompletedAt: string | null
  CreatedAt: string
  UpdatedAt: string
}

