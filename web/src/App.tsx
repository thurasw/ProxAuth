import useDarkMode from "use-dark-mode";
import Login from "./screens/Login";
import { useAuth } from "./hooks/AuthContext";
import Home from "./screens/Home";
import LoadingScreen from "./screens/Loading";

export default function App() {
    useDarkMode(true, {
        global: window,
        classNameDark: 'dark',
        classNameLight: 'light'
    });

    const { isAuthenticated, isLoading } = useAuth();

    if (isLoading) {
        return <LoadingScreen />
    }
    return (
        <div>
            <Login isOpen={!isAuthenticated} />
            { isAuthenticated && <Home /> }
        </div>
    );
}
