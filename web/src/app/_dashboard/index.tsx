import dynamic from 'next/dynamic';

export default dynamic(() => import('./DashboardWrap'), { ssr: false });