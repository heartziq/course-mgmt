import React, { useEffect, useState } from 'react';
import {useLocation, Redirect} from 'react-router-dom';
import qs from 'query-string';


const baseURL = "/api/v1/courses"

export default function CourseDetails({useAuth}) {
    let {token} = useAuth()
    const [course, setCourse] = useState({})
    const [status, setStatus] = useState(null)

    // Get APIKey from state
    let location = useLocation();
    // const {token} = location.state;
    // console.log("[CourseDetails] auth.token: ", token)

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
                    console.log("[CourseDetails] Status: ", httpResponse)
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

    // return status !== 403 ? <h1>CourseDetails: {course.Title}</h1>: (<Redirect
    //     to={{
    //         pathname: "/Dashboard/4",
    //         state: { message: "Please renew api key" }
    //     }}
    // />)
        
    //     return status !== 403 ? <h1>CourseDetails: {course.Title}</h1>: (<Redirect
    //     to={{
    //         pathname: "/Login",
    //         state: { message: "Please renew api key" }
    //     }}
    // />)

    let compo = renderDetails(status, course)
    return compo

// return status !== 403 ? <h1>CourseDetails: {course.Title}</h1>: <h1>403: renew or invalid</h1>
        

}

function renderDetails(statusCode, course) {
    let compo;
    switch (statusCode) {
        case 400:
            compo = <h1>Invalid API_KEY (null)</h1>
            break
        case 403:
            compo = <h1>Expired Key</h1>
            break
        case 401:
            compo = <h1>No key provided</h1>
            break
        default:
            compo = <h1>CourseDetails: {course.Title}</h1>
    }

    return compo
}