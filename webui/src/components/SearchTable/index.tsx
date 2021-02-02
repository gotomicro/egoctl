import React, {
  useEffect,
  useState,
  forwardRef,
  useImperativeHandle,
  Ref,
  ReactNode,
} from 'react';
import { Table, Form, message } from 'antd';
import { TableProps, ColumnGroupType, ColumnType } from 'antd/es/table';
import { FormInstance } from 'antd/es/form';
import { PaginationProps } from 'antd/es/pagination';
import { useLocation, useHistory } from 'umi';
import { stringify } from 'qs';

import './index.less';

export declare type TableColumnType<T> = ColumnGroupType<T> | ColumnType<T>;

export declare type TableEnumType = {
  text: ReactNode;
};

export declare type TableColumnTypes<T> = (TableColumnType<T> & {
  valueEnum?: {
    [key: string]: TableEnumType | ReactNode;
  };
})[];

function isTableEnumType(object: any): object is TableEnumType {
  return 'text' in object;
}

interface SearchTableProps<T>
  extends Omit<TableProps<T>, 'pagination' | 'dataSource' | 'columns'> {
  formContent: (
    search: (fields: any) => void,
    form: FormInstance,
  ) => React.ReactNode;
  request: (params: any) => Promise<{ data: T[]; pagination: PaginationProps }>;
  queryArrayFormat?: 'indices' | 'brackets' | 'repeat' | 'comma';
  columns: TableColumnTypes<T>;
}

export interface SearchTableInstance {
  refresh: (params?: any) => void;
  form: FormInstance;
}

const SearchTable = <T extends {}>(
  props: SearchTableProps<T>,
  ref: Ref<SearchTableInstance>,
) => {
  const {
    columns,
    formContent,
    request,
    queryArrayFormat,
    ...restProps
  } = props;
  const location = useLocation();
  const [dataSource, setDataSource] = useState<T[]>();
  const [pagination, setPagination] = useState<PaginationProps>();
  const [loading, setLoading] = useState<boolean>(true);
  //@ts-ignore
  const params = location.query;
  const [fields, setFields] = useState<any>({ ...params });
  const history = useHistory();
  const [form] = Form.useForm();

  useImperativeHandle(ref, () => ({
    refresh: params => {
      setFields({ ...fields, ...params });
    },
    form,
  }));

  const onSearch = (fields: any) => {
    setFields({ ...fields });
  };

  // 支持表格字段枚举显示
  for (let i = 0; i < columns.length; ++i) {
    let col = columns[i];
    if (!col.render && col.valueEnum) {
      col.render = val => {
        let title = val
        title = col.valueEnum.map((item, idx) => {
          if (val == item.value) {
            return item.title
          }
        })
        return title
        // console.log("val",val)
        // console.log("valueEnum",col.valueEnum)
        //
        // let enumVal = (col.valueEnum || {})[val];
        // console.log("enumVal",enumVal)
        //
        // if (!enumVal) return val;
        //
        // console.log("enumVal",enumVal)
        //
        // if (isTableEnumType(enumVal)) {
        //   return enumVal.title;
        // }
        //
        // return enumVal;
      };
    }
  }

  useEffect(() => {
    history.push({
      search:
        '?' +
        stringify(
          {
            ...params,
            ...fields,
          },
          { arrayFormat: queryArrayFormat || 'repeat' },
        ),
    });

    setLoading(true);

    request({ ...fields })
      .then(r => {
        setDataSource(r.data);
        setPagination(r.pagination);
        setLoading(false);
      })
      .catch(e => {
        message.error(e);
      });
    form.setFieldsValue({ ...fields });
  }, [fields]);

  return (
    <div>
      <div>{formContent(onSearch, form)}</div>

      <Table<T>
        style={{ marginTop: '10px' }}
        {...restProps}
        loading={loading}
        columns={columns}
        dataSource={dataSource}
        pagination={{
          ...pagination,
          onChange: (page, pageSize) => {
            setFields({
              ...fields,
              current: page,
              pageSize,
            });
          },
        }}
      />
    </div>
  );
};

export default forwardRef(SearchTable);
