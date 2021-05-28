import React, { useEffect, useState } from 'react';
import {useLocation, Redirect} from 'react-router-dom';
import qs from 'query-string';


const baseURL = "http://localhost:5000/api/v1/courses"

export default function CourseDetails() {
 
    const [course, setCourse] = useState({})
    const [status, setStatus] = useState(null)

    // Get APIKey from state
    let location = useLocation();
    const {token} = location.state;

    // Get courseId from query param
    let {id} = qs.parse(location.search)

    useEffect(
        () => {
            async function fetchData() {

                try {
                    let httpResponse = await fetch(
                        baseURL + `/${id}?key=${token}`,
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