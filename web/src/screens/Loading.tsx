import { Card, Spinner } from "@nextui-org/react";

export default function LoadingScreen() {
    return (
        <div className='w-screen h-screen flex items-center justify-center'>
            <Card
                isBlurred
                className='px-28 py-12 backdrop-saturate-100'
            >
                <Spinner color='white' size='lg' />
                <span className='font-mono mt-5'>Please wait...</span>
            </Card>
        </div>
    )
}
