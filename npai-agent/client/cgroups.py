#!/usr/bin/env python
# -*- coding: utf-8 -*-
# __author__ = "Wilson Lan"

import os
import multiprocessing


#####################
# utility functions #
#####################
def must_make_cgroup(hierarchy, cgroup_name):
    # create cgroup under specific hierarchy
    cgroup_tpl = "/sys/fs/cgroup/{hierarchy}/{cgroup_name}"
    cgroup = cgroup_tpl.format(hierarchy=hierarchy,
                               cgroup_name=cgroup_name)
    if not os.path.exists(cgroup):
        os.mkdir(cgroup, 0755)
    return cgroup

def set_cgroups_procs(pid, *cgroups):
    # add task to specific cgroup
    for cg in cgroups:
        procs = "{cgroup}/cgroup.procs".format(cgroup=cg)
        with open(procs, "a") as fp:
            fp.write("{pid}\n".format(pid=pid))

def set_cgroups_cfs_attr(*cgroups, **attrs):
    # attributes' values
    period_us = attrs.get("cfs_period_us", 100000)
    quota_us  = attrs.get("cfs_quota_us", -1)

    for cg in cgroups:
        hierarchy = os.path.basename(cg[:cg.rfind("/")])
        subsystem = hierarchy
        if subsystem == "cpu":
            cfs_quota_us  = "{cgroup}/cpu.cfs_quota_us".format(cgroup=cg)
            cfs_period_us = "{cgroup}/cpu.cfs_period_us".format(cgroup=cg)
            with open(cfs_quota_us, "w") as fp:
                fp.write("{quota_us}\n".format(quota_us=quota_us))
            with open(cfs_period_us, "w") as fp:
                fp.write("{period_us}\n".format(period_us=period_us))

def set_cgroups_cpuset_attr(*cgroups, **attrs):
    # attributes' values
    core_id = attrs.get("core_id", "0-{cpu_num}"
            .format(cpu_num=multiprocessing.cpu_count()-1))
    mem_node = attrs.get("mem_node", "0")

    for cg in cgroups:
        hierarchy = os.path.basename(cg[:cg.rfind("/")])
        subsystem = hierarchy
        if subsystem == "cpuset":
            cpus = "{cgroup}/cpuset.cpus".format(cgroup=cg)
            with open(cpus, "w") as fp:
                fp.write("{cpus}\n".format(cpus=core_id))
            mems = "{cgroup}/cpuset.mems".format(cgroup=cg)
            with open(mems, "w") as fp:
                fp.write("{mems}\n".format(mems=mem_node))

def clean_cgroups(*cgroups):
    for cg in cgroups:
        parent_cg = cg[:cg.rfind("/")]
        parent_procs = "{cgroup}/cgroup.procs".format(cgroup=parent_cg)
        with open(parent_procs, "a") as pfp:
            children_procs = "{cgroup}/cgroup.procs".format(cgroup=cg)
            with open(children_procs, "r") as cfp:
                procs = cfp.read()
            pfp.write("{procs}\n".format(procs=procs))
        os.rmdir(cg)


if __name__ == "__main__":
    pid = os.getpid()
    print "Task {pid}".format(pid=pid)

    # init related cgroup
    cgroups = []
    cgroups.append(must_make_cgroup("cpuset", "cgrouptest"))
    cgroups.append(must_make_cgroup("cpu", "cgrouptest"))

    # set cpuset cpus attributes
    set_cgroups_cpuset_attr(*cgroups, core_id="8")

    # set cgroup
    set_cgroups_procs(pid, *cgroups)

    # clean cgroup
    clean_cgroups(*cgroups)
