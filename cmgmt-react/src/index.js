import React, { createContext, useState, useContext, useEffect } from 'react';
import ReactDOM from 'react-dom';
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link
} from "react-router-dom";

import { CookiesProvider, useCookies } from 'react-cookie';


import CourseDetails from './components/CourseDetails';
import AllCourses from './components/AllCourses';
import Login from './components/Login';
import PrivateRoute from './components/PrivateRoute';
import Dashboard from './components/Dashboard';

// Page not found
import PageNotFound from './components/NotFound';


export default function Root() {
  return (
    <CookiesProvider>
      <App />
    </CookiesProvider>
  )
}
const authContext = createContext();

function Draft() {
  let auth = useAuth()
  useEffect(() =>
    fetch("/draft")
      .then(res => res.json())
      .then(({ token }) => {
        
        auth.signout(token)
      }), [])
  return (
    <div>
      <h1>coo</h1>
      <Link to="/Dashboard/4">
        to dashboard
      </Link>
    </div>

  )
}

function ProvideAuth({ children }) {
  const auth = useProvideAuth();
  return <authContext.Provider value={auth}>{children}</authContext.Provider>;
}
// const auth = useProvideAuth();
function useProvideAuth() {
  const [token, setToken] = useState(null);

  const signin = (value, callback) => {
    setToken(value)
    callback()
  };

  const signout = (token="") => {
    let newTokenValue = token || null
    setToken(newTokenValue)
  };

  return {
    token,
    signin,
    signout
  };
}

function useAuth() {
  return useContext(authContext);
}

function App() {
  const [cookies, setCookie] = useCookies(['name']);

  // function onChange(newName) {
  //   setCookie('name', newName, { path: '/' });
  // }
  useEffect(
    () => {
      setCookie("name", "HELLOOOO", { path: '/Login' })
    }
  )
  console.log("cookies: ", cookies.name)
  return (
    <ProvideAuth>
      <Router>
        <div>
          <nav>
            <ul>
              <li>
                <Link to="/">Browse All</Link>
              </li>
              <li>
                <Link to="/Login">Login</Link>
              </li>
              <li>
                <Link to="/draft">Get cookie</Link>
              </li>
            </ul>
          </nav>

          {/* A <Switch> looks through its children <Route>s and
          renders the first one that matches the current URL. */}
          <Switch>
            <Route path="/" exact >
              <AllCourses useAuth={useAuth} />
            </Route>
            <Route path="/Login">
              <Login useAuth={useAuth} />
            </Route>
            <PrivateRoute path="/CourseDetails" useAuth={useAuth}>
              <CourseDetails useAuth={useAuth} />
            </PrivateRoute>
            <Route path="/Dashboard/:id">
              <Dashboard useAuth={useAuth} />
            </Route>
            <Route path="/draft">
              <Draft />
            </Route>
            <Route component={PageNotFound} />
          </Switch>
        </div>
      </Router>
    </ProvideAuth>
  )
}
// ========================================
ReactDOM.render(
  <Root />,
  document.getElementById('root')
);
