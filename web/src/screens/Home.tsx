import { Dropdown, DropdownTrigger, Button, DropdownMenu, DropdownItem } from "@nextui-org/react";
import { IoAppsOutline, IoLogOut } from "react-icons/io5";
import { FaUserAstronaut } from "react-icons/fa";
import { useAuth } from "../hooks/AuthContext";
import LoadingScreen from "./Loading";
import { useState } from "react";

export default function Home() {

    const { user } = useAuth();
    const [ isLoggingOut, setIsLoggingOut ] = useState(false);
    
    if (isLoggingOut || !user) return <LoadingScreen />;
    return (
        <div className='h-screen max-h-screen w-screen overflow-hidden flex flex-col'>
            <UserDropdown setIsLoggingOut={setIsLoggingOut} />
            <div className='flex-grow'>

            </div>
        </div>
    )
}

function UserDropdown({ setIsLoggingOut }: { setIsLoggingOut: (b: boolean) => void }) {
    const { user, refreshUser } = useAuth();

    const logOut = async() => {
        await fetch('/api/auth/logout')
        refreshUser();
        setIsLoggingOut(true);
    }
    
    return (
        <div className='flex-shrink-0 flex justify-end px-5 py-3'>
            <Dropdown backdrop="blur">
                <DropdownTrigger>
                    <Button
                        variant="light"
                        className='bg-black/25 backdrop-blur-lg font-semibold font-mono text-lg'
                        startContent={<FaUserAstronaut size={18} />}
                    >
                        <span className='ms-1'>{ user!.username }</span>
                    </Button>
                </DropdownTrigger>
                <DropdownMenu variant="faded" aria-label="Static Actions">
                    <DropdownItem
                        key="edit-user"
                        startContent={<FaUserAstronaut size={18} />}
                    >
                        Edit User
                    </DropdownItem>
                    <DropdownItem
                        key="edit-apps"
                        startContent={<IoAppsOutline size={18} />}
                    >
                        Edit Dashboard
                    </DropdownItem>
                    <DropdownItem
                        key="signout"
                        className="text-danger"
                        color="danger"
                        startContent={<IoLogOut size={19} />}
                        onPress={logOut}
                    >
                        Sign Out
                    </DropdownItem>
                </DropdownMenu>
            </Dropdown>
        </div>
    )
}