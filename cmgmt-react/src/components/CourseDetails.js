import React, { useEffect, useState } from 'react';
import {useLocation} from 'react-router-dom';
import qs from 'query-string';


const baseURL = "http://localhost:5000/api/v1/courses"
// http://localhost:5000/api/v1/courses/LT5042?key=1469bb5bc0b9129857f01cb1ea8d7d1c6125e4cba38874a8b54af9e32a525351
export default function CourseDetails({ id }) {
    // console.log('id: ', id)
    // console.log('match', match.params)
    let location = useLocation();
    const {token} = location.state;
    // console.log('location: ', location)
    let params = qs.parse(location.search)
    console.log('CourseDetails, qs params: ', params.id)
    // console.log('token: ', token)
    const [course, setCourse] = useState({})
    // const [hasChanged, setHasChanged] = useState(false)

    useEffect(
        () => {
            async function fetchData() {
                let httpResponse = await fetch(
                    baseURL + `/${params.id}?key=${token}`,
                    { mode: "cors" }
                )
                let res = await httpResponse.json()

                setCourse(res)
            }
            fetchData();
            // setHasChanged(true)
            // console.log(course)
        },
        // () => fetch(
        //     baseURL + `/${id}?key=1469bb5bc0b9129857f01cb1ea8d7d1c6125e4cba38874a8b54af9e32a525351`,
        //     { mode: "cors" }
        // ).then(res => res.json())
        //     .then(result => setCourse(result))
        // ,
        [] // to prevent infinite loop
    )
    console.log(course)
    return (
        <h1>CourseDetails: {course.Title}</h1>

    )
}