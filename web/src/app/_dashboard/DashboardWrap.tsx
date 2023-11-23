"use client";
import { useColloState } from '@/hooks/useColloState';
import { NetworkGraph } from './NetworkGraph';
import { DateInput, FormComps, KeywordInput, Label, LoadingButton, StartButton, WrapProps } from './Forms';
import { FormEvent, useState } from 'react';

export default function Home() {
    const colloState = useColloState();
    const [loading, setLoading] = useState(false);
    const formwrapPorps: WrapProps = {
        onSubmit: function (event: FormEvent<HTMLFormElement>): void {
            event.preventDefault();
            setLoading(true);

            const form = new FormData(event.currentTarget);
            const start = form.get("start");
            const end = form.get("end");
            const keyword = form.get("keyword");
            if (!start || !end || !keyword) {
                setLoading(false)
                return;
            }

            // start search
            colloState.requestCollo(new Date(start.toString()), new Date(end.toString()), keyword.toString())
                .then(() => setLoading(false));
        }
    }
    return <>
        <FormComps.Wrap {...formwrapPorps}>
            <FormComps.Col>
                <Label htmlFor='keyword'>{"開始日"}</Label><DateInput id='start' name='start' defaultValue={"2023-03-01"} />
            </FormComps.Col>
            <FormComps.Col>
                <Label htmlFor='keyword'>{"終了日"}</Label><DateInput id='end' name='end' defaultValue={"2023-03-31"} />
            </FormComps.Col>
            <FormComps.Col>
                <Label htmlFor='keyword'>{"キーワード"}</Label><KeywordInput id='keyword' name='keyword' />
            </FormComps.Col>
            {loading ? <LoadingButton /> : <StartButton />}
        </FormComps.Wrap>
        <NetworkGraph {...colloState} />
    </>
}
