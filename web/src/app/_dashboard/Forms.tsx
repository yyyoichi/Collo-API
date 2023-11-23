"use client";
import React from 'react';

export type WrapProps = NonNullablePick<React.ComponentProps<"form">, "onSubmit">
    & Pick<React.ComponentProps<"form">, "children">;
const Wrap = ({ children, ...props }: WrapProps) => <form {...props}
    className='max-w-md mx-auto mt-8 p-6 bg-white shadow-md rounded-md'>{children}</form>
const Col = (props: Pick<React.ComponentProps<"div">, "children">) => <div
    className='mb-4'
>{props.children}</div>
export const FormComps = {
    Wrap, Col,
}

export type LabelProps = NonNullablePick<React.ComponentProps<"label">, "htmlFor">
    & Pick<React.ComponentProps<"form">, "children">;
export const Label = ({ children, ...props }: LabelProps) => <label {...props}
    className='block text-sm font-medium text-gray-600'
>{children}</label>

export type KeywordInputProps = NonNullablePick<React.ComponentProps<"input">, "id" | "name">
export const KeywordInput = (props: KeywordInputProps) => <input
    {...props}
    type='text'
    required
    className='mt-1 p-2 border rounded-md w-full focus:outline-none focus:border-blue-500'
/>

export type DateInputProps = NonNullablePick<React.ComponentProps<"input">, "id" | "name" | "defaultValue">
export const DateInput = (props: DateInputProps) => <input
    {...props}
    type='date'
    required
    className='mt-1 p-2 border rounded-md w-full focus:outline-none focus:border-blue-500'
/>

export const StartButton = () => <input
    type={"submit"}
    value={"開始"}
    className='block mx-auto bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline-blue cursor-pointer'
/>

export const LoadingButton = () => <button
    type='button'
    className="block mx-auto bg-gray-500 w-16 h-10 text-white font-bold rounded focus:outline-none focus:shadow-outline-blue disabled:opacity-50 disabled:cursor-not-allowed cursor-wait"
>
    <div className='h-4 w-4 m-auto animate-spin rounded-full border-b-2 border-t-2 cursor-wait' />
</button>