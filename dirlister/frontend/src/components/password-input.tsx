
import { Input } from "../components/ui/input";
import { Button } from "../components/ui/button";
import { useState } from "react";
import { Eye, EyeOff } from "lucide-react";

interface PasswordInputProps {
  placeholder: string;
  value: string;                                   
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void; 
}

export function PasswordInput({ placeholder, value, onChange }: PasswordInputProps) {
    const [visible, setVisible] =
        useState(false);

    return (
        <div className="flex gap-2">

            <Input
                placeholder={placeholder}
                type={
                    visible
                        ? "text"
                        : "password"
                }
                value={value}
                onChange={onChange}
                required

            />

            <Button
                type="button"
                variant="outline"
                onClick={() =>
                    setVisible(!visible)
                }
            >
                {visible
                    ? <Eye />
                    : <EyeOff />}
            </Button>

        </div>
    );
}