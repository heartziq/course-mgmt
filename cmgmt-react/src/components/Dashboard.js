import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';

export default function Dashboard({useAuth}) {
    let { id } = useParams()
    // let {token} = useAuth()
    // save cookie at state
    const [status, setStatus] = useState();
    useEffect(
        () => fetch("/dashboard/" + id, {
            // headers: {
            //     "Authorization": "Bearer " + token,
            // }
        })
            .then(
                res => setStatus(res.status),
                err => console.error(err)
            ),
        []
    )
    // load dashboard page here (if success)
    return status == 202 ? <h1>{id}, Renew APIKey button</h1> : <h1>login Failed</h1>
}