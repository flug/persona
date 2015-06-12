<?php

namespace Persona\Repository;


interface RepositoryInterface
{
    public function __construct($path, $profileName);

    public function processBuilder();

    public function getOutput();
}