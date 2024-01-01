"use client";
import React from 'react';

export type WrapProps = NonNullablePick<React.ComponentProps<"form">, "onSubmit">
    & Pick<React.ComponentProps<"form">, "children">;

export const FormComps = {
    Wrap: ({ children, ...props }: WrapProps) => <form {...props}
        className='absolute z-10 resize w-80 overflow-hidden max-w-md mx-auto p-6 bg-white shadow-md rounded-md'>{children}</form>,
    Col: (props: Pick<React.ComponentProps<"div">, "children">) => <div
        className='mb-4'
    >{props.children}</div>,
}

export type LabelProps = NonNullablePick<React.ComponentProps<"label">, "htmlFor">
    & Pick<React.ComponentProps<"form">, "children">;
export const Label = ({ children, ...props }: LabelProps) => <label {...props}
    className='block text-sm font-medium text-gray-600'
>{children}</label>

export type KeywordInputProps = NonNullablePick<React.ComponentProps<"input">, "id" | "name" | "defaultValue">
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

export const CheckboxLabel = ({ children, ...props }: LabelProps) => <label {...props}
    className="flex items-center text-xs font-medium text-gray-600"
>{children}</label>
export type PoSpeechCheckboxProps = NonNullablePick<React.ComponentProps<"input">, "id" | "name" | "value"> &
    Pick<React.ComponentProps<"input">, "defaultChecked">
export const PoSpeechCheckbox = (props: PoSpeechCheckboxProps) => <input
    {...props}
    className='mr-[.04rem]'
    type='checkbox'
/>

export type StopwordsTextareaProps = NonNullablePick<React.ComponentProps<"textarea">, "id" | "name" | "placeholder">
export const StopwordsTextarea = (props: StopwordsTextareaProps) => <textarea
    {...props}
    className='mt-1 p-2 border rounded-md text-sm w-full focus:outline-none focus:border-blue-500 resize-none'
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

type AccordionHeadProps = NonNullablePick<React.ComponentProps<"label">, "children">
type AccordionContentProps = NonNullablePick<React.ComponentProps<"div">, "children">
export const AccordionPanel = {
    Head: (props: AccordionHeadProps) => (
        <>
            <input type="checkbox" id="accordion" className='peer hidden' defaultChecked={true} />
            <label htmlFor='accordion' className='block w-auto text-right text-xs font-medium text-gray-600'>
                {props.children}
            </label>
        </>
    ),
    Content: (props: AccordionContentProps) => (
        <div className='peer-checked:h-0 peer-checked:opacity-0 peer-checked:invisible h-52 transition duration-500'>
            {props.children}
        </div>
    )
}