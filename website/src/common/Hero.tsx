/*
Copyright 2024 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
import { ReactNode } from "react";

export type HeroProps = {
  title: string;
  description?: ReactNode;
  children?: ReactNode;
};

export default function Hero({ title, description, children }: HeroProps) {
  return (
    <section className="flex flex-col items-center p-12">
      <div className="flex flex-col items-center gap-4 max-w-screen-lg">
        <h2 className="text-4xl md:text-6xl font-semibold text-primary mb-4">
          {title}
        </h2>
        <p className="md:my-6 leading-loose text-foreground/80">
          {description}
        </p>
        {children}
      </div>
    </section>
  );
}
