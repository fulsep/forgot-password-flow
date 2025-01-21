import React from 'react'
import { useForm } from 'react-hook-form'
import { Link, useNavigate, useSearchParams } from 'react-router-dom'

function Register() {
    const { register, handleSubmit } = useForm()
    const [sp] = useSearchParams()
    const navigate = useNavigate()
    const [msg, setMsg] = React.useState('')

    const token = sp.get("token")

    React.useEffect(() => {
        if (!token) {
            navigate("/forgot-password")
        }
    }, [token])
    const submitForm = (val) => {
        const body = new URLSearchParams({
            ...val,
            token
        })

        fetch("http://localhost:8888/reset-password", {
            method: "POST",
            body
        })
            .then(res => res.json())
            .then(res => {
                setMsg(res.message)
            })
    }
    return (
        <div className='flex h-screen justify-center items-center'>
            <form onSubmit={handleSubmit(submitForm)} className='max-w-md px-5 w-full flex flex-col gap-5'>
                {msg && <div role="alert" className="alert alert-warning">
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        className="h-6 w-6 shrink-0 stroke-current"
                        fill="none"
                        viewBox="0 0 24 24">
                        <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth="2"
                            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                    </svg>
                    <span>{msg}</span>
                </div>}
                <div>
                    <label className="form-control w-full">
                        <div className="label">
                            <span className="label-text">New Password</span>
                        </div>
                        <input type="password" placeholder="Type here" className="input input-bordered w-full" {...register("password")} />
                    </label>
                </div>
                <div>
                    <label className="form-control w-full">
                        <div className="label">
                            <span className="label-text">Confirm Password</span>
                        </div>
                        <input type="password" placeholder="Type here" className="input input-bordered w-full" {...register("confirm-password")} />
                    </label>
                </div>
                <div>
                    <button className="btn btn-primary w-full">Submit</button>
                </div>
                <div className='text-center'>
                    <Link to="/login">Already have an account? Login here</Link>
                </div>
            </form>
        </div>
    )
}

export default Register