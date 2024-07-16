import { Modal, ModalBody, ModalContent, ModalFooter, ModalHeader } from "@nextui-org/react";

export default function UserModal({ isOpen } : { isOpen: boolean }) {
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
                    Users
                </ModalHeader>
                <ModalBody>

                </ModalBody>
                <ModalFooter>

                </ModalFooter>
            </ModalContent>
        </Modal>
    )
}