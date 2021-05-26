import React, { useEffect, useState } from 'react';
import {useLocation, Redirect} from 'react-router-dom';
import qs from 'query-string';


const baseURL = "http://localhost:5000/api/v1/courses"
// http://localhost:5000/api/v1/courses/LT5042?key=1469bb5bc0b9129857f01cb1ea8d7d1c6125e4cba38874a8b54af9e32a525351

function renderError() {
    return (
        <h1>API key expired. pls renew</h1>
    )
}
export default function CourseDetails() {

    const [course, setCourse] = useState({})
    const [status, setStatus] = useState(null)

    let location = useLocation();
    const {token} = location.state;

    let params = qs.parse(location.search)

    useEffect(
        () => {
            async function fetchData() {

                try {
                    let httpResponse = await fetch(
                        baseURL + `/${params.id}?key=${token}`,
                        { mode: "cors" }
                    )
                    setStatus(httpResponse.status)
                    setCourse(httpResponse.json())
                } catch(err){
                    
                    console.error("erraaaa: ", err)
                }

            }

            fetchData();

        },
        [] // to prevent infinite loop
    )

    return status !== 403 ? <h1>CourseDetails: {course.Title}</h1>: (<Redirect
        to={{
            pathname: "/Login",
            state: { message: "Please renew api key" }
        }}
    />)
        

}