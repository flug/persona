<?php

namespace Persona;

use Symfony\Component\Console\Command\Command as SymfonyCommand;


class Command extends SymfonyCommand
{

    /**
     * @param $key
     * @param int $flags
     * @return mixed|string
     */
    public function get($key, $flags = 0)
    {
        return $this->getKernel()->get($key, $flags);
    }

    /**
     * @return static of AppKernel
     */
    public function getKernel()
    {
        return \AppKernel::getInstance();
    }
}