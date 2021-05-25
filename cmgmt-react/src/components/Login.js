import React, { useState } from 'react';
// import Cookies from 'js-cookie';
// import {useCookies} from 'react-cookie';

import {useHistory, useLocation} from 'react-router-dom';


export default function Login({ useAuth }) {
  let history = useHistory();
  let location = useLocation();
  let auth = useAuth()
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const baseURL = "http://localhost:5000/login?NewKey=False";

  async function handleSubmit(e) {
    e.preventDefault()

    try {
      let res = await fetch(baseURL, {
        method: "POST",
        mode: "cors",
        body: JSON.stringify({
          "username": username,
          "password": password,
        })
      })
      let apiKeyDetails = await res.json()
      let token = apiKeyDetails["access_token"]
      // console.log(token)

      // Get previous state.from
      let { from } = location.state || { from: { pathname: "/" } };
      console.log('Login.js: ', location)

      // Redirect (back) to state.from
      auth.signin(token, () => history.replace(from))

    } catch (err) {
      console.log("error occured while fetching")
      console.error(err)
    }
  }

  return (
    <form onSubmit={handleSubmit}>
      <label>
        Name:
            <input type="text" value={username} onChange={e => setUsername(e.target.value.trim())} />
      </label>
      <label>
        Password:
            <input type="password" value={password} onChange={e => setPassword(e.target.value.trim())} />
      </label>
      <input type="submit" value="Submit" />
    </form>
  );
}


// export default function Login() {
//     return (
//         <Fragment>
//         <h1>LOGIINININ Details</h1>

//       </Fragment>
//     )
// }
