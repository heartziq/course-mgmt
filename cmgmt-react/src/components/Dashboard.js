import React from 'react';
import {useParams} from 'react-router-dom';

export default function Dashboard(){
    let {id} = useParams()
    console.log(id)
    return <h1>{id}</h1>
}