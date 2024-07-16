import { ReactNode, useCallback, useMemo, useState } from "react";
import useSWR, { SWRConfig } from "swr";
import { User } from "../types";
import { AuthContext } from "../hooks/AuthContext";

export function AuthProvider({ children }: { children: ReactNode }) {
    /**
     * If the request errors out with a 401, we want to disable the fetcher, until the user logs in.
     */
    const [ fetchEnabled, setFetchEnabled ] = useState(true);
    const userFetcher = async (url: URL) => {
        const response = await fetch(url)
        if (response.status === 401) {
            setFetchEnabled(false);
            throw new Error("Unauthorized!")
        }

        const res = await response.json();
        return res;
    }

    // User state
    const { data, error, mutate } = useSWR<User>(
        fetchEnabled ? "/api/users" : null,
        userFetcher,
        {
            refreshInterval: 1000 * 60,
            shouldRetryOnError: false,
        }
    );

    // Reenalbe fetcher & manually set data as stale to refetch
    const refreshUser = useCallback(() => {
        setFetchEnabled(true);
        mutate();
    }, [ mutate ]);

    // Fetcher for all other queries
    const apiFetcher = async (rs: RequestInfo | URL, init?: RequestInit) => {
        const response = await fetch(rs, init);
        if (response.status === 401) {
            refreshUser();
            throw new Error("Unauthorized!")
        }

        const res = await response.json();
        return res;
    }

    const auth = useMemo(() => ({
        isAuthenticated: data !== undefined && error === undefined,
        user: error !== undefined ? null : (data || null),
        error: error,
        isLoading: fetchEnabled && !data && !error,
        refreshUser
    }), [data, error, refreshUser, fetchEnabled])

    return (
        <AuthContext.Provider value={auth}>
            <SWRConfig value={{ fetcher: apiFetcher }}>
                {children}
            </SWRConfig>
        </AuthContext.Provider>
    )
}
