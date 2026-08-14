export interface ApiFieldError {
  field: string
  code: string
  message: string
}

export interface ApiErrorDetails {
  code: string
  message: string
  request_id?: string
  fields?: ApiFieldError[]
}

export interface ApiErrorResponse {
  course_errors: ApiErrorDetails
}
