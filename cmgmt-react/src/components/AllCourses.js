import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';

const baseURL = "http://localhost:5000/api/v1/courses"

export default function AllCourses({ useAuth }) {
    // define states
    let auth = useAuth();

    const [error, setError] = useState(null);
    const [isLoaded, setIsLoaded] = useState(false);
    const [courseList, setCourseList] = useState([]);

    useEffect(
        () => fetch(baseURL).then(res => {
            console.log(res.status)
            return res.json()
        }).then(
            result => {

                setCourseList(result)
                setIsLoaded(true)
            },
            err => setError(err)
        ),
        []
    )

    if (error) {
        return <div>Error: {error.message}</div>
    } else if (!isLoaded) {
        return <div>Loading</div>
    } else {
        console.log(courseList);
        return (
            <ul>
                {
                    // How to pass props to Link.. {id and key}
                    courseList
                        .map(course => (
                            <li key={course.id}>
                                <Link
                                    to={{
                                        pathname: "/CourseDetails",
                                        search: `?id=${course.id}`,
                                        state: {
                                            fromDashboard: true,
                                            token: auth.token,
                                        }
                                    }}
                                >
                                    {course.Title}

                                </Link>
                            </li>


                        ))
                }
            </ul>
        )

    }
}