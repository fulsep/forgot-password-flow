import React from 'react'
import { Link } from 'react-router-dom'

function Home() {
    return (
        <div className='flex flex-col h-screen'>
            <div className="navbar bg-base-100 shadow">
                <div className="flex-1">
                    <Link to="/" className="btn btn-ghost text-xl">Forgot Password Demo</Link>
                </div>
                <div className="flex-none">
                    <ul className="menu menu-horizontal px-1">
                        <li><Link to="/login">Login</Link></li>
                        <li><Link to="/register">Register</Link></li>
                        <li><Link to="/forgot-password">Forgot Password</Link></li>
                    </ul>
                </div>
            </div>
            <div className='flex-1 flex-col gap-5 flex justify-center items-center text-3xl font-bold'>
                <div>Welcome!</div>
                <div>Check link on the Navigation Bar</div>
            </div>
        </div>
    )
}

export default Home