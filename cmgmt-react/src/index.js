import React, { createContext, useState, useContext } from 'react';
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

import { useEffect } from 'react';


export default function Root() {
  return (
    <CookiesProvider>
      <App />
    </CookiesProvider>
  )
}
const authContext = createContext();

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

  const signout = () => {
    setToken(null)
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
              <CourseDetails />
            </PrivateRoute>
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
