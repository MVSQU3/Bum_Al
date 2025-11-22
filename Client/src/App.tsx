import { useEffect } from "react";
import { Navigate, Route, Routes } from "react-router-dom";
import { useAuthStore } from "./store/userAuthStore";
import Login from "./Login";
import Home from "./Home";
import AddAlbum from "./AddAlbum";
import UpdateAlbum from "./UpdateAlbum";

function App() {
  const { authUser, checkAuth } = useAuthStore();

  useEffect(() => {
    checkAuth();
  }, []);
  
  console.log("authUser = ", authUser);

  return (
    <>
      <Routes>
        <Route
          path="/"
          element={authUser ? <Home /> : <Navigate to={"/login"} />}
        />
        <Route
          path="/album/add"
          element={authUser ? <AddAlbum /> : <Login />}
        />
        <Route
          path="/album/update"
          element={authUser ? <UpdateAlbum /> : <Login />}
        />
        <Route
          path="/login"
          element={!authUser ? <Login /> : <Navigate to={"/"} />}
        />
      </Routes>
    </>
  );
}

export default App;
