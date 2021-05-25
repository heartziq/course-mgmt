import React from 'react';
import { Redirect, Route } from 'react-router-dom';

export default function PrivateRoute({ useAuth, children, ...rest }) {
    // subscribe to 'auth' context
    let auth = useAuth()

    return (
        <Route
            token={auth.token}
            {...rest}
            render={({ location }) => {
                // console.log('PrivateRoute.js', location)
                if (auth.token)
                    return (children)
                else
                    return (<Redirect
                        to={{
                            pathname: "/Login",
                            state: { from: location }
                        }}
                    />)
                }
            }

        />
    )
}