import { Button, Input, Modal, ModalBody, ModalContent, ModalFooter, ModalHeader } from "@nextui-org/react";
import { useState } from "react";
import { useAuth } from "../hooks/AuthContext";
import { IoEye, IoEyeOff, IoLockClosed } from "react-icons/io5";
import { TiUser } from "react-icons/ti";
import { RiErrorWarningFill } from "react-icons/ri";

export default function Login({ isOpen }: { isOpen: boolean }) {

    const { refreshUser } = useAuth();
    const [ isVisible, setIsVisible ] = useState(false);
    const [ isLoading, setIsLoading ] = useState(false);
    const [ err, setErr ] = useState('');

    const [ username, setUsername ] = useState('');
    const [ password, setPassword ] = useState('');

    const handleLogin = async() => {
        setIsLoading(true)
        try {
            const res = await fetch("/api/auth/login", {
                headers: {
                    'content-type': 'application/json'
                },
                method: 'POST',
                body: JSON.stringify({ username, password })
            });
            if (!res.ok) {
                throw new Error(await res.text());
            }

            refreshUser();
        }
        catch(err: any) {
            console.error(err);
            setErr(err.message);
        }
        finally {
            setIsLoading(false);
            setPassword("");
        }
    };

    return (
        <Modal
            isOpen={isOpen}
            isDismissable={false}
            isKeyboardDismissDisabled={true}
            backdrop='blur'
            hideCloseButton={true}
        >
            <ModalContent>
                <ModalHeader className="flex flex-col gap-1 py-5 text-2xl text-center">
                    Sign in to continue
                </ModalHeader>
                <form
                    onSubmit={(e) => {
                        e.preventDefault();
                        handleLogin();
                    }}
                >
                    <ModalBody>
                        <Input
                            value={username}
                            onChange={(e) => setUsername(e.target.value)}
                            label="Username"
                            variant="bordered"
                            labelPlacement="outside"
                            type="text"
                            startContent={<TiUser className="text-2xl text-default-400 pointer-events-none" />}
                        />
                        <Input
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            label="Password"
                            variant="bordered"
                            labelPlacement="outside"
                            type={isVisible ? "text" : "password"}
                            startContent={<IoLockClosed className="text-xl text-default-400 pointer-events-none" />}
                            endContent={
                                <button className="focus:outline-none" type="button" onClick={() => setIsVisible(t => !t)}>
                                    {isVisible ? (
                                        <IoEyeOff className="text-2xl text-default-400 pointer-events-none" />
                                    ) : (
                                        <IoEye className="text-2xl text-default-400 pointer-events-none" />
                                    )}
                                </button>
                            }
                        />
                        {
                            err !== '' && (
                                <div className='flex items-center mt-3 px-3 py-1.5 w-full rounded-medium bg-danger/20 dark:text-danger-500'>
                                    <RiErrorWarningFill size={22} className='my-1' />
                                    <span className='ms-3 text-small text-wrap font-mono'>{ err }</span>
                                </div>
                            )
                        }
                    </ModalBody>
                    <ModalFooter>
                        <Button
                            color="primary"
                            variant="shadow"
                            fullWidth={true}
                            isLoading={isLoading}
                            type='submit'
                        >
                            Login
                        </Button>
                    </ModalFooter>
                </form>
            </ModalContent>
        </Modal>
    )
}