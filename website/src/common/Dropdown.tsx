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

import React, { useEffect, useRef, useState } from "react";
import useClickOutside from "../hooks/useClickOutside";
import Icon from "../common/Icon";
import { twMerge } from "tailwind-merge";
import { Button } from "@/components/ui/button";

interface ContainerProps {
  children: React.ReactNode;
  defaultIndex?: number;
  onChange: (selected: { value: string }) => void;
  className?: string;
  placeholder?: string;
}

interface OptionProps {
  children: React.ReactNode;
  className?: string;
  onClick?: () => void;
  isSelected?: boolean;
}

interface OptionElementProps extends OptionProps {
  index: number;
}

function Container({
  children,
  defaultIndex = 0,
  onChange,
  className,
  placeholder = "Select an option",
}: ContainerProps) {

  const [open, setOpen] = useState(false);
  const [selectedIndex, setSelectedIndex] = useState(defaultIndex || 0);

  const ref = useRef<HTMLDivElement | null>(null);

  let items = React.Children.map(children, (child) => `${child}`) || [];

  useClickOutside(ref as React.MutableRefObject<HTMLElement>, () =>
    setOpen(false)
  );

  function handleOptionClick(index: number) {
    setOpen(false);
    setSelectedIndex(index);
  }

  useEffect(() => {
    selectedIndex < 0 && setSelectedIndex(0);
  }, [selectedIndex]);

  useEffect(() => {
    onChange({ value: items[selectedIndex] });
  }, [selectedIndex]);

  return (
    <div ref={ref} className="relative flex gap-2">
      <Button
        variant="outline"
        onClick={() => setOpen(!open)}
        className={twMerge(
          "flex items-center h-full justify-center duration-inherit hover:bg-transparent",
          className
        )}
      >
        {`${items[selectedIndex] || placeholder}`}
        <Icon
          icon="expand_more"
          className={twMerge(
            "pt-[3px] ml-2 scale-150 duration-300",
            open && "rotate-180"
          )}
        />
      </Button>

      <>
        <div
          className={twMerge(
            "flex flex-col left-1/2 -translate-x-1/2 z-10 duration-300 absolute top-full"
          )}
          style={{
            clipPath: !open
              ? "polygon(0% 0%, 0% 0%, 100% 0%, 100% 0%)"
              : "polygon(0% 0%, 0% 100%, 100% 100%, 100% 0%)",
          }}
        >
          {React.Children.map(children, (child, index) => {
            if (React.isValidElement(child) && child.type === Option)
              return React.cloneElement(child, {
                index,
                onClick: () => handleOptionClick(index),
                isSelected: index === selectedIndex,
              } as OptionElementProps);
          })}
        </div>
      </>
    </div>
  );
}

function Option({ children, className, onClick, isSelected }: OptionProps) {
  return (
    <Button
      className={twMerge(
        "dark:text-white text-black",
        className,
        isSelected && "bg-primary"
      )}
      onClick={onClick}
    >
      {children}
    </Button>
  );
}

const Dropdown = { Container, Option };

export default Dropdown;
