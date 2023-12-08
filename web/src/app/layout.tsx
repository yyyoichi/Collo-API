import '../styles/globals.css'
import type { Metadata } from 'next'
import { Inter } from 'next/font/google'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'Create Next App',
  description: 'Generated by create next app',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="jp">
      <body className={inter.className}>
        <main className='overflow-x-hidden w-full bg-white'>
          {children}
        </main>
        <footer className='flex flex-col items-center justify-center h-40 bg-slate-700'>
          <div className='text-yellow-50 my-4'>
            国会議事録をネットワーク分析してみる
          </div>
          <div className='text-gray-500 text-sm'>
            <a href='https://kokkai.ndl.go.jp' className='border-b border-gray-500'>国会議事録API</a>を使用しています。
          </div>
          <div className='text-gray-500  text-sm'>
            辞書として<a href='https://clrd.ninjal.ac.jp' className='border-b border-gray-500'>UniDic</a>を使用しています。
          </div>
        </footer>
      </body>
    </html>
  )
}
