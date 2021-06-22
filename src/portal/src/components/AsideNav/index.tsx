import React, { useEffect, useState } from 'react';
import { Menu, Icon } from 'cess-ui';
import { useHistory } from 'react-router-dom';
import { connect } from 'dva';

import './index.less';
import { ICommonState, useSelector } from 'umi';
const { SubMenu } = Menu;

interface NavItem {
  label: string;
  key: string;
  route: string;
  icon?: string;
  childrens?: NavItem[];
}

function getMenus(isMobile: boolean): NavItem[] {
  // 移动端只展示应用商店

  return [
    {
      label: '开发环境',
      key: 'dev_env',
      route: `/devEnvironment`,
      icon: 'lm-environment',
    },
    {
      label: '应用管理',
      key: 'application',
      route: 'application',
      icon: 'lm-applications',
      childrens: [
        {
          label: '应用市场',
          key: 'application_chart',
          route: '/application_chart',
        },
        {
          label: '应用实例',
          key: 'application_instance',
          route: '/application_instance',
        },
      ],
    },
    {
      label: '文件管理',
      key: 'file',
      route: '/file',
      icon: 'lm-file',
    },
    {
      label: '镜像管理',
      key: 'mirror',
      route: 'mirror',
      icon: 'lm-mirror',
      childrens: [
        {
          label: '私有镜像仓库',
          key: 'private_mirror',
          route: '/private_mirror',
        },
        {
          label: '公有镜像仓库',
          key: 'public_mirror',
          route: '/public_mirror',
        },
      ],
    },
  ];
}

function AsideNav() {
  const { isMobile }: ICommonState = useSelector((state: any) => state.common);
  const history = useHistory();
  const [selectedKeys, setSelectedKeys] = useState<string[]>([]);
  const [menus, setMenus] = useState<NavItem[]>([]);

  function handleClick(param: any) {
    setSelectedKeys([param.key]);
    history.push(param.key);
  }

  useEffect(() => {
    const key = history.location.pathname.split('/').slice(0, 2).join('/');
    if (key === '/') {
      // 路径为空的重定向
      history.push('/home');
    } else {
      setSelectedKeys([key]);
    }
  }, [history.location]);
  useEffect(() => {
    setMenus(getMenus(isMobile));
  }, [isMobile]);
  return (
    <aside>
      <Menu
        onClick={(param) => handleClick(param)}
        defaultSelectedKeys={['home']}
        defaultOpenKeys={['mirror', 'application']}
        mode="inline"
        selectedKeys={selectedKeys}
      >
        {menus.map((m) => {
          return m.childrens ? (
            <SubMenu
              key={m.route}
              title={
                <span>
                  <Icon type={m.icon} />
                  <span>{m.label}</span>
                </span>
              }
            >
              {m.childrens &&
                m.childrens.map((c) => <Menu.Item key={c.route}>{c.label}</Menu.Item>)}
            </SubMenu>
          ) : (
            <Menu.Item key={m.route}>
              {' '}
              <span>
                <Icon type={m.icon} />
                <span>{m.label}</span>
              </span>
            </Menu.Item>
          );
        })}
      </Menu>
    </aside>
  );
}

const mapStateTopProps = (state: any) => {
  return {
    ...state,
  };
};

export default connect(mapStateTopProps)(AsideNav);