import { lazy, Suspense, type ComponentType } from 'react'
import { createBrowserRouter, Navigate } from 'react-router-dom'
import { ProtectedRoute } from '../features/auth/ui/ProtectedRoute'
import { AppLayout } from '../shared/ui/AppLayout'
import { Spinner } from '../shared/ui/Spinner'

const LoginPage = lazy(() => import('../pages/login/LoginPage').then((module) => ({ default: module.LoginPage })))
const RegisterPage = lazy(() => import('../pages/register/RegisterPage').then((module) => ({ default: module.RegisterPage })))
const DashboardPage = lazy(() => import('../pages/dashboard/DashboardPage').then((module) => ({ default: module.DashboardPage })))
const CoursesPage = lazy(() => import('../pages/courses/CoursesPage').then((module) => ({ default: module.CoursesPage })))
const CourseDetailsPage = lazy(() => import('../pages/course-details/CourseDetailsPage').then((module) => ({ default: module.CourseDetailsPage })))
const MyCoursesPage = lazy(() => import('../pages/my-courses/MyCoursesPage').then((module) => ({ default: module.MyCoursesPage })))
const LessonPage = lazy(() => import('../pages/lesson/LessonPage').then((module) => ({ default: module.LessonPage })))
const CreatorCoursesPage = lazy(() => import('../pages/creator/CreatorCoursesPage').then((module) => ({ default: module.CreatorCoursesPage })))
const CreateCoursePage = lazy(() => import('../pages/creator/CreateCoursePage').then((module) => ({ default: module.CreateCoursePage })))
const ManageCoursePage = lazy(() => import('../pages/creator/ManageCoursePage').then((module) => ({ default: module.ManageCoursePage })))
const NotFoundPage = lazy(() => import('../pages/not-found/NotFoundPage').then((module) => ({ default: module.NotFoundPage })))

function routeElement(Page: ComponentType) {
  return <Suspense fallback={<main className="route-loading"><Spinner /><p>Загружаем страницу…</p></main>}><Page /></Suspense>
}

export const router = createBrowserRouter([
  { path: '/login', element: routeElement(LoginPage) },
  { path: '/register', element: routeElement(RegisterPage) },
  {
    element: <AppLayout />,
    children: [
      { index: true, element: <Navigate to="/courses" replace /> },
      { path: '/courses', element: routeElement(CoursesPage) },
      { path: '/courses/:courseId', element: routeElement(CourseDetailsPage) },
      { path: '/courses/:courseId/lessons/:lessonId', element: routeElement(LessonPage) },
      { element: <ProtectedRoute />, children: [
        { path: '/dashboard', element: routeElement(DashboardPage) },
        { path: '/my-courses', element: routeElement(MyCoursesPage) },
        { path: '/creator', element: routeElement(CreatorCoursesPage) },
        { path: '/creator/courses/new', element: routeElement(CreateCoursePage) },
        { path: '/creator/courses/:courseId', element: routeElement(ManageCoursePage) },
      ] },
      { path: '*', element: routeElement(NotFoundPage) },
    ],
  },
])
