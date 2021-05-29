import React, { Fragment, useState } from 'react';

import { useHistory, useLocation, Link } from 'react-router-dom';


export default function Login({ useAuth }) {
  let history = useHistory();
  let location = useLocation();
  let auth = useAuth()
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  let { message } = location.state || ""


  const baseURL = "/login?NewKey=False";

  async function handleSubmit(e) {
    e.preventDefault()

    try {
      let res = await fetch(baseURL, {
        method: "POST",
        body: JSON.stringify({
          "username": username,
          "password": password,
        })
      })
      let apiKeyDetails = await res.json()
      let token = apiKeyDetails["access_token"]
      console.log("token: ", token)

      // Get previous state.from
      let { from } = location.state || { from: { pathname: "/" } };

      // Update auth.token and,
      // Redirect (back) to state.from
      auth.signin(token, () => history.replace(from))
      // auth.signout(token)

    } catch (err) {
      console.log("error occured while fetching")
      console.error(err)
    }
  }

  return (<Fragment>
    <p>{message}</p>
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
  </Fragment>)


}