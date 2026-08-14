# LMS demonstration regression checklist

## Authentication

- Register a student and confirm redirect to the dashboard.
- Log out, log in again, and refresh a protected page.
- Confirm an expired access token refreshes once and the original request retries.
- Confirm logged-out access to `/dashboard`, `/my-courses`, and `/creator` redirects to login.

## Learner flow

- Open `/courses` logged out; verify loading, empty/error handling, and course cards.
- Open a course and confirm the author fallback contains no UUID.
- Log in as a student, enroll, and confirm the course appears in `/my-courses` without refresh.
- Open lessons in order; verify previous/next navigation and plain-text content.
- Complete a lesson; verify progress refreshes without a full page reload.
- Complete the final lesson; verify 100% progress and completed-course state.

## Creator flow

- Create a draft course and edit title, description, and level.
- Add, edit, reposition, and delete draft lessons.
- Refresh the draft management URL; confirm course and lessons persist.
- Confirm publishing is disabled with zero lessons and works with at least one.
- Refresh a published course; confirm lessons are visible and structurally read-only.
- Confirm published metadata remains editable.
- Archive the course; confirm it leaves the public catalog but remains in creator management.
- Confirm a closed course and its lessons are read-only.

## Routing and responsive checks

- Directly open and refresh `/courses/:id`, `/courses/:id/lessons/:lessonId`, `/my-courses`, and `/creator/courses/:id` through nginx.
- Check header navigation, forms, cards, lesson navigation, and dialogs at desktop and mobile widths.
- Confirm logout clears protected server-state views and returns to public navigation.

