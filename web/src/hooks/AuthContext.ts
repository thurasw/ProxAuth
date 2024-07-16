import { createContext, useContext } from "react";
import { User } from "../types";

export interface AuthContextProps {
    isAuthenticated: boolean;
    isLoading: boolean;
    user: User | null;
    refreshUser: () => void;
}
export const AuthContext = createContext<AuthContextProps | undefined>(undefined);

export function useAuth() {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error('useAuth must be used within a AuthProvider');
    }
    return context;

}