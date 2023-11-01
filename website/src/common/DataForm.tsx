/*
Copyright 2023 The Vitess Authors.

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

import React, { useEffect, useState } from "react";

interface ContainerProps {
  children?: React.ReactNode;
  onSubmit?: (data: object) => void;
}

//Form implementation by @marsian83 (https://github.com/marsian83)
function Container(
  props: ContainerProps & React.FormHTMLAttributes<HTMLFormElement>
) {
  const [data, setData] = useState<object>({});

  function submitHandler(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    props.onSubmit && props.onSubmit(data);
  }

  return (
    <form {...props} onSubmit={submitHandler}>
      <MappedInputs setData={setData}>{props.children}</MappedInputs>
    </form>
  );
}

function MappedInputs(props: {
  children: React.ReactNode;
  setData: React.Dispatch<React.SetStateAction<object>>;
}) {
  useEffect(() => {
    React.Children.forEach(props.children, (child) => {
      if (
        React.isValidElement(child) &&
        child.type === Input &&
        child.props.type !== "submit" &&
        child.props.name
      ) {
        props.setData((p) => ({
          ...p,
          [child.props.name]:
            child.props.defaultValue || child.props.value || "",
        }));
      }
    });
  }, []);

  return (
    <>
      {React.Children.map(props.children, (child, key) => {
        if (React.isValidElement(child))
          if (child.type === Input)
            return React.cloneElement(child as React.ReactElement, {
              key,
              onChange: (event: React.ChangeEvent<HTMLInputElement>) => {
                props.setData((p) => ({
                  ...p,
                  [child.props.name]: event.target.value,
                }));
              },
            });
          else if (child.props.children)
            return React.cloneElement(child as React.ReactElement, {
              key,
              children: (
                <MappedInputs setData={props.setData}>
                  {child.props.children}
                </MappedInputs>
              ),
            });
        return child;
      })}
    </>
  );
}

interface InputProps {
  name: string;
}

type InputAttributes = React.InputHTMLAttributes<HTMLInputElement>;

function Input(
  props: InputProps &
    Partial<InputAttributes> &
    Required<Pick<InputAttributes, "name">>
) {
  return <input {...props} />;
}

const DataForm = { Container, Input };

export default DataForm;
